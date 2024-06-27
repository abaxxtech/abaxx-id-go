package dwn

import (
	"io"

	"github.com/abaxxtech/abaxx-id-go/internal/types"
)

type HandlerRequest struct {
	Tenant     string
	Message    types.GenericMessage
	DataStream io.Reader
}
