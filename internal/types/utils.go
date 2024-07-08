package types

type GenericMessage struct {
	Descriptor struct {
		Interface        string `json:"interface"`
		Method           string `json:"method"`
		MessageTimestamp string `json:"messageTimestamp"`
	} `json:"descriptor"`
	Authorization struct {
		Signer    string `json:"signer"`
		Signature string `json:"signature"`
	} `json:"authorization"`
}

type GenericSignaturePayload struct {
	DescriptorCid      string `json:"descriptorCid"`
	DelegatedGrantId   string `json:"delegatedGrantId,omitempty"`
	PermissionsGrantId string `json:"permissionsGrantId,omitempty"`
	ProtocolRole       string `json:"protocolRole,omitempty"`
}
