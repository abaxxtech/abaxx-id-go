package dwn

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	jwk "github.com/lestrrat-go/jwx/v2/jwk"
)

const authSig = `
{ "signature": 
  {
    "payload": "eyJkZXNjcmlwdG9yQ2lkIjoiYmFmeXJlaWJpc3F2eG9sNHF6ZDVkem9wa3VlZHR3eTRrbnp1Z2dudWprbW15b25hNzV3dGxvd2k0ZGkifQ",
    "signatures": [
      {
	"protected": "eyJraWQiOiJkaWQ6aW9uOkVpQVRUNmZLQWc0cHJzQVE5VkE5ZGcxbmZHSTBLdUtfVTNzUUdrcFZvVDlNdkE6ZXlKa1pXeDBZU0k2ZXlKd1lYUmphR1Z6SWpwYmV5SmhZM1JwYjI0aU9pSnlaWEJzWVdObElpd2laRzlqZFcxbGJuUWlPbnNpY0hWaWJHbGpTMlY1Y3lJNlczc2lhV1FpT2lKa2QyNHRjMmxuSWl3aWNIVmliR2xqUzJWNVNuZHJJanA3SW1OeWRpSTZJa1ZrTWpVMU1Ua2lMQ0pyZEhraU9pSlBTMUFpTENKNElqb2lkMGRwVW5WdlRsVkZjbkZvU0d0WFoweGpiSFJHTW5sb1pYWjNhMGxTTTNaQk9FTkNTSGswTkhWU2F5SjlMQ0p3ZFhKd2IzTmxjeUk2V3lKaGRYUm9aVzUwYVdOaGRHbHZiaUpkTENKMGVYQmxJam9pU25OdmJsZGxZa3RsZVRJd01qQWlmU3g3SW1sa0lqb2laSGR1TFdWdVl5SXNJbkIxWW14cFkwdGxlVXAzYXlJNmV5SmpjbllpT2lKelpXTndNalUyYXpFaUxDSnJkSGtpT2lKRlF5SXNJbmdpT2lKYVdHeEZUWGMzTkVoUlRYUllWamRxZGs1aFoyNUZSR294WlRKUWVVSXpkVGgyYVdsUlRXaFRZbG8wSWl3aWVTSTZJbVZJVDI1NlYyYzNVRVZsWjFaUFVuRk5OakZDZDFVeGRsRTBMVzlyTVdOcVFWcDFja3AzUkhneU9VVWlmU3dpY0hWeWNHOXpaWE1pT2xzaWEyVjVRV2R5WldWdFpXNTBJbDBzSW5SNWNHVWlPaUpLYzI5dVYyVmlTMlY1TWpBeU1DSjlYU3dpYzJWeWRtbGpaWE1pT2x0N0ltbGtJam9pWkhkdUlpd2ljMlZ5ZG1salpVVnVaSEJ2YVc1MElqcDdJbVZ1WTNKNWNIUnBiMjVMWlhseklqcGJJaU5rZDI0dFpXNWpJbDBzSW01dlpHVnpJanBiSW1oMGRIQnpPaTh2WkhkdU1DNWtaWFl1WVdKaGVIZ3VhV1FpTENKb2RIUndjem92TDJSM2JqRXVaR1YyTG1GaVlYaDRMbWxrSWl3aWFIUjBjSE02THk5a2QyNHlMbVJsZGk1aFltRjRlQzVwWkNKZExDSnphV2R1YVc1blMyVjVjeUk2V3lJalpIZHVMWE5wWnlKZGZTd2lkSGx3WlNJNklrUmxZMlZ1ZEhKaGJHbDZaV1JYWldKT2IyUmxJbjFkZlgxZExDSjFjR1JoZEdWRGIyMXRhWFJ0Wlc1MElqb2lSV2xDWTJneE4ySmZZbWN6TVhwelpIbDJkMVZtYlRreU1GWnVVbXBCU3pOa1JUaHZUR2RuYVRZd0xUaENRU0o5TENKemRXWm1hWGhFWVhSaElqcDdJbVJsYkhSaFNHRnphQ0k2SWtWcFF6RTFjbFJCVTBWWVNuaHFhMHBrU1Y5MGFtTlhhRTltWkhoU1RIZDVNSGg2ZDFsdE1XSmFVa2xTYldjaUxDSnlaV052ZG1WeWVVTnZiVzFwZEcxbGJuUWlPaUpGYVVOb016RnBVUzFaU1dJd1FYTkVibGRZYkVzM2NVZEVaWFJXVVU1aFMwMTJSM0F3Y1dwYVJtVlFjVEJuSW4xOSNkd24tc2lnIiwiYWxnIjoiRWREU0EifQ",
	"signature": "fqWuqWMEpdv2cf0KvozKasOmEd-QJLNlvvuMqNzVvJSbrrTx2EzKIr1IesfBoxdnmdeRkFivODsGVMuNVwjlBQ"
      }
    ]
  }
}
`


func NewTestDwn() *Dwn {
	dwn, err := NewDwn(DwnConfig{
		DidResolver:  nil,
		TenantGate:   NewAllowAllTenantGate(),
		MessageStore: NewMemoryMessageStore(),
		DataStore:    NewMemoryDatastore(),
		EventLog:     NewMemoryEventLog(),
	})

	if err != nil {
		panic(err)
	}
	return dwn
}

func TestReadDid(t *testing.T) {
	content, err := os.ReadFile("testdata/did.json")
	assert.NoError(t, err)

	didInfo := fileDID{}
	err = json.Unmarshal(content, &didInfo)
	assert.NoError(t, err)

	assert.Equal(t, "ryan.rawson@abaxx.tech", didInfo.Id)
	assert.Equal(t, "ryan.rawson@abaxx.tech", didInfo.Metadata.Email)

	// find a key:
	// re-serialize a key to import it:
	printKey(t, didInfo, "#dwn-enc")
	printKey(t, didInfo, "#dwn-sig")
	
	
}

func printKey(t *testing.T, didInfo fileDID, keyid string) {
	keystruct, err := didInfo.GetKeyById(keyid)
	assert.NoError(t, err)

	t.Log("keystruct", keystruct)

	keyBytes, err := json.Marshal(keystruct)
	assert.NoError(t, err)

	k1, err := jwk.ParseKey(keyBytes)
	assert.NoError(t, err)
	t.Log("jwk k1", k1)

	t.Log("keyid", k1.KeyID())
	t.Log("keyops", k1.KeyOps())
	t.Log("keytype", k1.KeyType())
	t.Log("algo", k1.Algorithm().String())

	x, err := json.MarshalIndent(k1, "", "  ")

	assert.NoError(t, err)
	t.Log("serialize", string(x))
}

func TestCreateAuthObjects(t *testing.T) {
	auth := PlainAuthorization{}
	err := json.Unmarshal([]byte(authSig), &auth)
	if err != nil {
		t.Error(err)
	}

	genSig := GeneralJws{Payload: "lkjasdlfjk",
		Signatures: []Signature{{"abc", "def"}}}

	stubDwn := NewTestDwn()
	err = stubDwn.authenticate(auth)
	t.Log("Error from authenticate", err)
	if err != nil {
		t.Error(err)
	}

	//t.Log("auth", auth)

	_, err = json.Marshal(auth)
	if err != nil {
		t.Error("json marshal of auth failed", err)
	}
	//t.Log("json", string(r))

	auth2 := AuthorizationDelegatedGrant{
		genSig, &DelegatedGrant{
			PlainAuthorization{genSig},
			DelegatedGrantDescriptor{"recordid", "encodeddata"}}}
	_, err = json.Marshal(auth2)

	//t.Log("auth2", string(r))
}
