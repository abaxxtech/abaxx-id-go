package store

import (
	"context"
	"errors"
	"io"

	importer "github.com/ipfs/boxo/ipld/unixfs/importer"
	blockservice "github.com/ipfs/go-blockservice"
	cid "github.com/ipfs/go-cid"
	ds "github.com/ipfs/go-datastore"
	nsds "github.com/ipfs/go-datastore/namespace"
	dsquery "github.com/ipfs/go-datastore/query"
	dsleveldb "github.com/ipfs/go-ds-leveldb"
	blockstore "github.com/ipfs/go-ipfs-blockstore"
	chunker "github.com/ipfs/go-ipfs-chunker"
	offline "github.com/ipfs/go-ipfs-exchange-offline"
	files "github.com/ipfs/go-ipfs-files"
	format "github.com/ipfs/go-ipld-format"
	dag "github.com/ipfs/go-merkledag"
	unixfs "github.com/ipfs/go-unixfs"
	uio "github.com/ipfs/go-unixfs/io"
)

// PlaceholderValue is used as a placeholder value for reference counting
var PlaceholderValue = []byte{0}

// DataStoreLevelConfig holds configuration options for DataStoreLevel
type DataStoreLevelConfig struct {
	BlockstoreLocation string
}

// DataStoreLevel is a simple implementation of DataStore that works in both the browser and server-side.
// Leverages LevelDB under the hood.
//
// It has the following structure (`+` represents a sublevel and `->` represents a key->value pair):
//
//	'data' + <tenant> + <dataCid> -> <data>
//	'references' + <tenant> + <dataCid> + <messageCid> -> PlaceholderValue
//
// This allows for the <data> to be shared for everything that uses the same <dataCid> while also making
// sure that the <data> can only be deleted if there are no <messageCid> for any <tenant> still using it.
type DataStoreLevel struct {
	config    DataStoreLevelConfig
	datastore ds.Batching
}

// NewDataStoreLevel creates a new instance of DataStoreLevel
func NewDataStoreLevel(config DataStoreLevelConfig) (*DataStoreLevel, error) {
	if config.BlockstoreLocation == "" {
		config.BlockstoreLocation = "data/DATASTORE"
	}

	// Initialize LevelDB datastore
	levelDB, err := dsleveldb.NewDatastore(config.BlockstoreLocation, nil)
	if err != nil {
		return nil, err
	}

	return &DataStoreLevel{
		config:    config,
		datastore: levelDB,
	}, nil
}

// Open opens the datastore (no-op for LevelDB)
func (d *DataStoreLevel) Open() error {
	// LevelDB datastore does not require opening
	return nil
}

// Close closes the datastore
func (d *DataStoreLevel) Close() error {
	return d.datastore.Close()
}

// Put stores the data and updates reference counting
func (d *DataStoreLevel) Put(ctx context.Context, tenant, messageCid, dataCid string, dataReader io.Reader) (*PutResult, error) {
	refDS := d.getDatastoreForReferenceCounting(tenant, dataCid)
	dataBS := d.getBlockstoreForStoringData(tenant, dataCid)

	// Add reference
	refKey := ds.NewKey(messageCid)
	if err := refDS.Put(ctx, refKey, PlaceholderValue); err != nil {
		return nil, err
	}

	// Import data into blockstore
	dagService := dag.NewDAGService(blockservice.New(dataBS, offline.Exchange(dataBS)))
	file := files.NewReaderFile(dataReader)

	// params := helpers.DagBuilderParams{
	// 	Dagserv:    dagService,
	// 	RawLeaves:  true,
	// 	CidBuilder: cid.V1Builder{Codec: cid.DagProtobuf, MhType: multihash.SHA2_256},
	// }

	// dagBuilder, _ := params.New(chunker.NewSizeSplitter(file, int64(1024*256))) // 256KB chunks

	rootNode, err := importer.BuildDagFromReader(dagService, chunker.DefaultSplitter(file))
	if err != nil {
		return nil, err
	}

	dataSize, err := sizeOfNode(rootNode)
	if err != nil {
		return nil, err
	}

	return &PutResult{
		DataCid:  rootNode.Cid().String(),
		DataSize: dataSize,
	}, nil
}

// Get retrieves the data if the caller has access
func (d *DataStoreLevel) Get(ctx context.Context, tenant, messageCid, dataCid string) (*GetResult, error) {
	refDS := d.getDatastoreForReferenceCounting(tenant, dataCid)
	dataBS := d.getBlockstoreForStoringData(tenant, dataCid)

	// Check if messageCid is allowed
	refKey := ds.NewKey(messageCid)
	hasRef, err := refDS.Has(ctx, refKey)
	if err != nil || !hasRef {
		return nil, errors.New("access denied or reference not found")
	}

	// Check if data exists
	c, err := cid.Decode(dataCid)
	if err != nil {
		return nil, err
	}

	hasData, err := dataBS.Has(ctx, c)
	if err != nil || !hasData {
		return nil, errors.New("data not found")
	}

	dagService := dag.NewDAGService(blockservice.New(dataBS, offline.Exchange(dataBS)))
	rootNode, err := dagService.Get(ctx, c)
	if err != nil {
		return nil, err
	}

	reader, err := uio.NewDagReader(ctx, rootNode, dagService)
	if err != nil {
		return nil, err
	}

	dataSize, err := sizeOfNode(rootNode)
	if err != nil {
		return nil, err
	}

	return &GetResult{
		DataCid:    rootNode.Cid().String(),
		DataSize:   dataSize,
		DataReader: reader,
	}, nil
}

// Delete removes the reference and deletes data if it's no longer referenced
func (d *DataStoreLevel) Delete(ctx context.Context, tenant, messageCid, dataCid string) error {
	refDS := d.getDatastoreForReferenceCounting(tenant, dataCid)
	dataBS := d.getBlockstoreForStoringData(tenant, dataCid)

	// Delete reference
	refKey := ds.NewKey(messageCid)
	if err := refDS.Delete(ctx, refKey); err != nil {
		return err
	}

	// Check if there are any remaining references
	keys, err := refDS.Query(ctx, dsquery.Query{})
	if err != nil {
		return err
	}
	defer keys.Close()

	if keys.Next() != nil {
		// References still exist, do not delete data
		return nil
	}

	// Delete data
	return dataBS.DeleteBlock(ctx, cid.MustParse(dataCid))
}

// Clear deletes everything in the store
func (d *DataStoreLevel) Clear(ctx context.Context) error {
	return d.datastore.Close() // Close and reopen to clear
}

// PutResult is the result of a Put operation
type PutResult struct {
	DataCid  string
	DataSize uint64
}

// GetResult is the result of a Get operation
type GetResult struct {
	DataCid    string
	DataSize   uint64
	DataReader io.ReadCloser
}

// Helper functions

// getDatastoreForReferenceCounting returns the datastore used for reference counting
func (d *DataStoreLevel) getDatastoreForReferenceCounting(tenant, dataCid string) ds.Datastore {
	referencesDS := nsds.Wrap(d.datastore, ds.NewKey("references"))
	tenantDS := nsds.Wrap(referencesDS, ds.NewKey(tenant))
	return nsds.Wrap(tenantDS, ds.NewKey(dataCid))
}

// getBlockstoreForStoringData returns the blockstore used for storing data
func (d *DataStoreLevel) getBlockstoreForStoringData(tenant, dataCid string) blockstore.Blockstore {
	dataDS := nsds.Wrap(d.datastore, ds.NewKey("data"))
	tenantDS := nsds.Wrap(dataDS, ds.NewKey(tenant))
	dataCidDS := nsds.Wrap(tenantDS, ds.NewKey(dataCid))
	return blockstore.NewBlockstore(dataCidDS)
}

// sizeOfNode calculates the total size of a DAG node
func sizeOfNode(node format.Node) (uint64, error) {
	if unixfsNode, ok := node.(*dag.ProtoNode); ok {
		fsNode, err := unixfs.FSNodeFromBytes(unixfsNode.Data())
		if err != nil {
			return 0, err
		}
		return fsNode.FileSize(), nil
	}
	return 0, errors.New("unsupported node type")
}
