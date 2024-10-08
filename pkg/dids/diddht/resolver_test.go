package diddht

import (
	"encoding/hex"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_resolve(t *testing.T) {
	// test vector
	vectors := map[string]string{
		"did:dht:9tjoow45ef1hksoo96bmzkwwy3mhme95d7fsi3ezjyjghmp75qyo": "ea33e704f3a48a3392f54b28744cdfb4e24780699f92ba7df62fd486d2a2cda3f263e1c6bcbd" +
			"75d438be7316e5d6e94b13e98151f599cfecefad0b37432bd90a0000000065b0ed1600008400" +
			"0000000300000000035f6b30045f6469643439746a6f6f773435656631686b736f6f3936626d" +
			"7a6b777779336d686d653935643766736933657a6a796a67686d70373571796f000010000100" +
			"001c2000373669643d303b743d303b6b3d5f464d49553174425a63566145502d437536715542" +
			"6c66466f5f73665332726c4630675362693239323445045f747970045f6469643439746a6f6f" +
			"773435656631686b736f6f3936626d7a6b777779336d686d653935643766736933657a6a796a" +
			"67686d70373571796f000010000100001c2000070669643d372c36045f6469643439746a6f6f" +
			"773435656631686b736f6f3936626d7a6b777779336d686d653935643766736933657a6a796a" +
			"67686d70373571796f000010000100001c20002726763d303b766d3d6b303b617574683d6b30" +
			"3b61736d3d6b303b64656c3d6b303b696e763d6b30",
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		did := "did:dht:" + r.URL.Path[1:]
		defer r.Body.Close()
		buf, ok := vectors[did]
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		data, err := hex.DecodeString(buf)
		assert.NoError(t, err)
		_, err = w.Write(data)
		assert.NoError(t, err)

	}))
	defer ts.Close()

	r := NewResolver(ts.URL, http.DefaultClient)

	for did := range vectors {
		t.Run(did, func(t *testing.T) {
			res, err := r.Resolve(did)
			assert.NoError(t, err)
			assert.NotZero(t, res.Document)
			assert.Equal(t, res.Document.ID, did)
		})
	}
}
