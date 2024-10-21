package jwk_test

import (
	"encoding/json"
	"testing"

	"github.com/abaxxtech/abaxx-id-go/pkg/jwk"
	"github.com/stretchr/testify/assert"
)

func TestJWK_MarshalJSON(t *testing.T) {
	key := jwk.JWK{
		KTY: "EC",
		CRV: "P-256",
		X:   "MKBCTNIcKUSDii11ySs3526iDZ8AiTo7Tu6KPAqv7D4",
		Y:   "4Etl6SRW2YiLUrN5vfvVHuhp7x8PxltmWWlbbM4IFyM",
	}

	b, err := json.Marshal(&key)
	assert.NoError(t, err)

	var unmarshalled map[string]interface{}
	err = json.Unmarshal(b, &unmarshalled)
	assert.NoError(t, err)

	assert.Equal(t, "EC", unmarshalled["kty"])
	assert.Equal(t, "P-256", unmarshalled["crv"])
	assert.Equal(t, "MKBCTNIcKUSDii11ySs3526iDZ8AiTo7Tu6KPAqv7D4", unmarshalled["x"])
	assert.Equal(t, "4Etl6SRW2YiLUrN5vfvVHuhp7x8PxltmWWlbbM4IFyM", unmarshalled["y"])
}

func TestJWK_UnmarshalJSON(t *testing.T) {
	jsonStr := `{
		"kty": "EC",
		"crv": "P-256",
		"x": "MKBCTNIcKUSDii11ySs3526iDZ8AiTo7Tu6KPAqv7D4",
		"y": "4Etl6SRW2YiLUrN5vfvVHuhp7x8PxltmWWlbbM4IFyM",
		"use": "sig",
		"kid": "1"
	}`

	var key jwk.JWK
	err := json.Unmarshal([]byte(jsonStr), &key)
	assert.NoError(t, err)

	assert.Equal(t, "EC", key.KTY)
	assert.Equal(t, "P-256", key.CRV)
	assert.Equal(t, "MKBCTNIcKUSDii11ySs3526iDZ8AiTo7Tu6KPAqv7D4", key.X)
	assert.Equal(t, "4Etl6SRW2YiLUrN5vfvVHuhp7x8PxltmWWlbbM4IFyM", key.Y)
}
