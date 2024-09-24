package utils

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/btcsuite/btcutil/base58"
)

type Convert struct {
	data   interface{}
	format string
}

func NewConvert(data interface{}, format string) *Convert {
	return &Convert{data: data, format: format}
}

func ArrayBuffer(data []byte) *Convert {
	return NewConvert(data, "ArrayBuffer")
}

func Base58Btc(data string) *Convert {
	return NewConvert(data, "Base58Btc")
}

func Base64Url(data string) *Convert {
	return NewConvert(data, "Base64Url")
}

func BufferSource(data []byte) *Convert {
	return NewConvert(data, "BufferSource")
}

func Hex(data string) *Convert {
	if len(data)%2 != 0 {
		panic("Hex input must have an even number of characters.")
	}
	return NewConvert(data, "Hex")
}

func Multibase(data string) *Convert {
	return NewConvert(data, "Multibase")
}

func Object(data map[string]interface{}) *Convert {
	return NewConvert(data, "Object")
}

func String(data string) *Convert {
	return NewConvert(data, "String")
}

func Uint8Array(data []byte) *Convert {
	return NewConvert(data, "Uint8Array")
}

func (c *Convert) ToArrayBuffer() []byte {
	switch c.format {
	case "Base58Btc":
		return base58.Decode(c.data.(string))
	case "Base64Url":
		decoded, _ := base64.RawURLEncoding.DecodeString(c.data.(string))
		return decoded
	case "BufferSource", "Uint8Array":
		return c.data.([]byte)
	case "Hex":
		decoded, _ := hex.DecodeString(c.data.(string))
		return decoded
	case "String":
		return []byte(c.data.(string))
	default:
		panic(fmt.Sprintf("Conversion from %s to ArrayBuffer is not supported.", c.format))
	}
}

func (c *Convert) ToBase58Btc() string {
	switch c.format {
	case "ArrayBuffer", "Uint8Array":
		return base58.Encode(c.data.([]byte))
	case "Multibase":
		return c.data.(string)[1:]
	default:
		panic(fmt.Sprintf("Conversion from %s to Base58Btc is not supported.", c.format))
	}
}

func (c *Convert) ToBase64Url() string {
	switch c.format {
	case "ArrayBuffer", "BufferSource", "Uint8Array":
		return base64.RawURLEncoding.EncodeToString(c.data.([]byte))
	case "Object":
		jsonData, _ := json.Marshal(c.data)
		return base64.RawURLEncoding.EncodeToString(jsonData)
	case "String":
		return base64.RawURLEncoding.EncodeToString([]byte(c.data.(string)))
	default:
		panic(fmt.Sprintf("Conversion from %s to Base64Url is not supported.", c.format))
	}
}

func (c *Convert) ToHex() string {
	switch c.format {
	case "ArrayBuffer", "Base64Url", "Uint8Array":
		return hex.EncodeToString(c.ToUint8Array())
	default:
		panic(fmt.Sprintf("Conversion from %s to Hex is not supported.", c.format))
	}
}

func (c *Convert) ToMultibase() string {
	switch c.format {
	case "Base58Btc":
		return "z" + c.data.(string)
	default:
		panic(fmt.Sprintf("Conversion from %s to Multibase is not supported.", c.format))
	}
}

func (c *Convert) ToObject() map[string]interface{} {
	switch c.format {
	case "Base64Url":
		decoded, _ := base64.RawURLEncoding.DecodeString(c.data.(string))
		var obj map[string]interface{}
		json.Unmarshal(decoded, &obj)
		return obj
	case "String":
		var obj map[string]interface{}
		json.Unmarshal([]byte(c.data.(string)), &obj)
		return obj
	case "Uint8Array":
		var obj map[string]interface{}
		json.Unmarshal(c.data.([]byte), &obj)
		return obj
	default:
		panic(fmt.Sprintf("Conversion from %s to Object is not supported.", c.format))
	}
}

func (c *Convert) ToString() string {
	switch c.format {
	case "ArrayBuffer", "Uint8Array":
		return string(c.data.([]byte))
	case "Base64Url":
		decoded, _ := base64.RawURLEncoding.DecodeString(c.data.(string))
		return string(decoded)
	case "Object":
		jsonData, _ := json.Marshal(c.data)
		return string(jsonData)
	default:
		panic(fmt.Sprintf("Conversion from %s to String is not supported.", c.format))
	}
}

func (c *Convert) ToUint8Array() []byte {
	switch c.format {
	case "ArrayBuffer", "BufferSource", "Uint8Array":
		return c.data.([]byte)
	case "Base58Btc":
		return base58.Decode(c.data.(string))
	case "Base64Url":
		decoded, _ := base64.RawURLEncoding.DecodeString(c.data.(string))
		return decoded
	case "Hex":
		decoded, _ := hex.DecodeString(c.data.(string))
		return decoded
	case "Object":
		jsonData, _ := json.Marshal(c.data)
		return jsonData
	case "String":
		return []byte(c.data.(string))
	default:
		panic(fmt.Sprintf("Conversion from %s to Uint8Array is not supported.", c.format))
	}
}

func RemoveUndefinedProperties(obj map[string]interface{}) {
	for key, value := range obj {
		if value == nil {
			delete(obj, key)
		} else if nestedObj, ok := value.(map[string]interface{}); ok {
			RemoveUndefinedProperties(nestedObj)
		}
	}
}
