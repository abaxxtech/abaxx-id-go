package store

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/ipfs/go-cid"
	ipld "github.com/ipfs/go-ipld-format"
	"github.com/ipfs/go-unixfs"
	"github.com/ipfs/go-unixfs/importer/balanced"
	"github.com/ipfs/go-unixfs/importer/helpers"
	"github.com/ipfs/go-unixfs/importer/trickle"
	"github.com/multiformats/go-multihash"
)

// PlaceholderValue is used as a placeholder value for reference counting
var PlaceholderValue = []byte{}

// DataStoreLevel implements the DataStore interface using LevelDB
type DataStoreLevel struct {
	blockstore *BlockstoreLevel
}

// NewDataStoreLevel creates a new DataStoreLevel
func NewDataStoreLevel(blockstoreLocation string) (*DataStoreLevel, error) {
	blockstore, err := NewBlockstoreLevel(blockstoreLocation)
	if err != nil {
		return nil, fmt.Errorf("failed to create blockstore: %w", err)
	}

	return &DataStoreLevel{
		blockstore: blockstore,
	}, nil
}

// Open opens the underlying blockstore
func (d *DataStoreLevel) Open() error {
	return d.blockstore.Open()
}

// Close closes the underlying blockstore
func (d *DataStoreLevel) Close() error {
	return d.blockstore.Close()
}

// Put stores data in the datastore
func (d *DataStoreLevel) Put(ctx context.Context, tenant, messageCid, dataCid string, dataStream io.Reader) (PutResult, error) {
	refBlockstore, err := d.getBlockstoreForReferenceCounting(tenant, dataCid)
	if err != nil {
		return PutResult{}, fmt.Errorf("failed to get reference blockstore: %w", err)
	}

	err = refBlockstore.Put(ctx, cid.NewCidV1(cid.Raw, multihash.IDENTITY), PlaceholderValue)
	if err != nil {
		return PutResult{}, fmt.Errorf("failed to put reference: %w", err)
	}

	dataBlockstore, err := d.getBlockstoreForStoringData(tenant, dataCid)
	if err != nil {
		return PutResult{}, fmt.Errorf("failed to get data blockstore: %w", err)
	}

	dagService := ipld.NewDAGService(dataBlockstore)
	params := unixfs.DagBuilderParams{
		Dagserv:    dagService,
		RawLeaves:  true,
		Maxlinks:   helpers.DefaultLinksPerBlock,
		NoCopy:     false,
		CidBuilder: cid.V1Builder{Codec: cid.DagProtobuf, MhType: multihash.SHA2_256},
	}

	db, err := params.New(io.NopCloser(dataStream))
	if err != nil {
		return PutResult{}, fmt.Errorf("failed to create DAG builder: %w", err)
	}

	var nd ipld.Node
	if params.Trickle {
		nd, err = trickle.Layout(db)
	} else {
		nd, err = balanced.Layout(db)
	}
	if err != nil {
		return PutResult{}, fmt.Errorf("failed to create DAG layout: %w", err)
	}

	size, err := nd.Size()
	if err != nil {
		return PutResult{}, fmt.Errorf("failed to get node size: %w", err)
	}

	return PutResult{
		DataCid:  nd.Cid().String(),
		DataSize: uint64(size),
	}, nil
}

// Get retrieves data from the datastore
func (d *DataStoreLevel) Get(ctx context.Context, tenant, messageCid, dataCid string) (GetResult, error) {
	refBlockstore, err := d.getBlockstoreForReferenceCounting(tenant, dataCid)
	if err != nil {
		return GetResult{}, err
	}

	allowed, err := refBlockstore.Has(ctx, cid.NewCidV1(cid.Raw, multihash.IDENTITY))
	if err != nil {
		return GetResult{}, err
	}
	if !allowed {
		return GetResult{}, errors.New("not allowed")
	}

	dataBlockstore, err := d.getBlockstoreForStoringData(tenant, dataCid)
	if err != nil {
		return GetResult{}, err
	}

	c, err := cid.Decode(dataCid)
	if err != nil {
		return GetResult{}, err
	}

	exists, err := dataBlockstore.Has(ctx, c)
	if err != nil {
		return GetResult{}, err
	}
	if !exists {
		return GetResult{}, errors.New("data not found")
	}

	dagService := ipld.NewDAGService(dataBlockstore)
	nd, err := dagService.Get(ctx, c)
	if err != nil {
		return GetResult{}, err
	}

	file, err := unixfs.ExtractDataFromNode(nd)
	if err != nil {
		return GetResult{}, err
	}

	size, err := nd.Size()
	if err != nil {
		return GetResult{}, err
	}

	return GetResult{
		DataCid:    c.String(),
		DataSize:   uint64(size),
		DataStream: file,
	}, nil
}

// Associate associates a message CID with a data CID
func (d *DataStoreLevel) Associate(ctx context.Context, tenant, messageCid, dataCid string) (AssociateResult, error) {
	refBlockstore, err := d.getBlockstoreForReferenceCounting(tenant, dataCid)
	if err != nil {
		return AssociateResult{}, err
	}

	isEmpty, err := refBlockstore.IsEmpty()
	if err != nil {
		return AssociateResult{}, err
	}
	if isEmpty {
		return AssociateResult{}, errors.New("no existing reference")
	}

	dataBlockstore, err := d.getBlockstoreForStoringData(tenant, dataCid)
	if err != nil {
		return AssociateResult{}, err
	}

	c, err := cid.Decode(dataCid)
	if err != nil {
		return AssociateResult{}, err
	}

	exists, err := dataBlockstore.Has(ctx, c)
	if err != nil {
		return AssociateResult{}, err
	}
	if !exists {
		return AssociateResult{}, errors.New("data not found")
	}

	err = refBlockstore.Put(ctx, cid.NewCidV1(cid.Raw, multihash.IDENTITY), PlaceholderValue)
	if err != nil {
		return AssociateResult{}, err
	}

	dagService := ipld.NewDAGService(dataBlockstore)
	nd, err := dagService.Get(ctx, c)
	if err != nil {
		return AssociateResult{}, err
	}

	size, err := nd.Size()
	if err != nil {
		return AssociateResult{}, err
	}

	return AssociateResult{
		DataCid:  c.String(),
		DataSize: uint64(size),
	}, nil
}

// Delete removes a message CID association and potentially the data if it's the last reference
func (d *DataStoreLevel) Delete(ctx context.Context, tenant, messageCid, dataCid string) error {
	refBlockstore, err := d.getBlockstoreForReferenceCounting(tenant, dataCid)
	if err != nil {
		return err
	}

	err = refBlockstore.Delete(ctx, cid.NewCidV1(cid.Raw, multihash.IDENTITY))
	if err != nil {
		return err
	}

	isEmpty, err := refBlockstore.IsEmpty()
	if err != nil {
		return err
	}
	if !isEmpty {
		return nil
	}

	dataBlockstore, err := d.getBlockstoreForStoringData(tenant, dataCid)
	if err != nil {
		return err
	}

	return dataBlockstore.Clear()
}

// Clear deletes everything in the store
func (d *DataStoreLevel) Clear() error {
	return d.blockstore.Clear()
}

func (d *DataStoreLevel) getBlockstoreForReferenceCounting(tenant, dataCid string) (*BlockstoreLevel, error) {
	refCountingPartition, err := d.blockstore.Partition("references")
	if err != nil {
		return nil, err
	}
	tenantPartition, err := refCountingPartition.Partition(tenant)
	if err != nil {
		return nil, err
	}
	return tenantPartition.Partition(dataCid)
}

func (d *DataStoreLevel) getBlockstoreForStoringData(tenant, dataCid string) (*BlockstoreLevel, error) {
	dataPartition, err := d.blockstore.Partition("data")
	if err != nil {
		return nil, err
	}
	tenantPartition, err := dataPartition.Partition(tenant)
	if err != nil {
		return nil, err
	}
	return tenantPartition.Partition(dataCid)
}

// PutResult represents the result of a Put operation
type PutResult struct {
	DataCid  string
	DataSize uint64
}

// GetResult represents the result of a Get operation
type GetResult struct {
	DataCid    string
	DataSize   uint64
	DataStream io.Reader
}

// AssociateResult represents the result of an Associate operation
type AssociateResult struct {
	DataCid  string
	DataSize uint64
}
