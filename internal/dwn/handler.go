package dwn

import (
	"errors"
	"io"

	"github.com/abaxxtech/abaxx-id-go/internal/types"
	cid "github.com/ipfs/go-cid"
	mh "github.com/multiformats/go-multihash"
)

type HandlerRequest struct {
	Tenant     string
	Message    types.GenericMessage
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
		// TODO verify these values and put them into a const
		if cid.Type() != 85 {
			return errors.New("Bad code")
		}
		if cid.Type() != 1 {
			return errors.New("Only support V1 for CIDs")
		}

		// TODO Verify these values and put them into a const
		de, err := mh.Decode(cid.Hash())
		if err != nil {
			return err
		}
		if de.Code != 18 {
			return errors.New("Only support sha2-256")
		}
		if de.Name != "sha2-256" {
			return errors.New("Only support sha2-256")
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
		// TODO serialize here
		results[i] = res
	}
	// TODO return results here.
	return nil
}
