package store

type GenericMessage interface{}

type DataFilter struct {
	Property string
	Operator string
	Value    interface{}
}

// type KeyValues map[string]interface{}
