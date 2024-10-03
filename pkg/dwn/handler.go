package dwn

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"

	cid "github.com/ipfs/go-cid"
	mh "github.com/multiformats/go-multihash"
)

const (
	expectedCidType    = 85
	expectedCidVersion = 1
	expectedHashCode   = 18
	expectedHashName   = "sha2-256"
)

type HandlerRequest struct {
	Tenant     string
	Message    map[string]interface{}
	DataStream io.Reader
}

func (message *RecordsWrite) Handle(dwn *Dwn) error {
	return nil
}

func validateCids(cids []string) error {

	for _, cidStr := range cids {
		cid, err := cid.Decode(cidStr)
		if err != nil {
			return err
		}

		if cid.Type() != expectedCidType {
			return errors.New("invalid CID type")
		}
		if cid.Version() != expectedCidVersion {
			return errors.New("only support v1 for CIDs")
		}

		de, err := mh.Decode(cid.Hash())
		if err != nil {
			return err
		}
		if de.Code != expectedHashCode {
			return errors.New("only support sha2-256")
		}
		if de.Name != expectedHashName {
			return errors.New("only support sha2-256")
		}
	}

	return nil
}

func (message *MessagesGet) Handle(tenant Tenant, dwn *Dwn) error {
	// Already done.
	// messagesGet = await MessagesGet.parse(message)

	// validateMessageCids
	err := validateCids(message.Descriptor.MessageCids)
	if err != nil {
		return err
	}
	err = dwn.authenticate(message.Authorization)
	if err != nil {
		return err
	}
	err = dwn.authorize(tenant, message.Authorization)
	if err != nil {
		return err
	}

	// set of message cids....
	results := make([]interface{}, len(message.Descriptor.MessageCids))
	for i, cid := range message.Descriptor.MessageCids {
		res, err := dwn.messageStore.Get(tenant, MessageCid(cid))
		if err != nil {
			return err
		}
		// Serialize the message to JSON
		serializedMessage, err := json.Marshal(res)
		if err != nil {
			return fmt.Errorf("failed to serialize message: %w", err)
		}
		results[i] = json.RawMessage(serializedMessage)
	}
	// TODO return results here.
	return nil
}
