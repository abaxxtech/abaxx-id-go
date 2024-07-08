# `crypto` <!-- omit in toc -->

# Table of Contents <!-- omit in toc -->

- [Features](#features)
- [Usage](#usage)
  - [`dsa`](#dsa)
    - [Key Generation](#key-generation)
    - [Signing](#signing)
    - [Verifying](#verifying)
- [Directory Structure](#directory-structure)
  - [Rationale](#rationale)

# Features

* secp256k1 keygen, deterministic signing, and verification
* ed25519 keygen, signing, and verification
* higher-level API for `ecdsa` (Elliptic Curve Digital Signature Algorithm)
* higher-level API for `eddsa` (Edwards-Curve Digital Signature Algorithm) 
* higher level API for `dsa` in general (Digital Signature Algorithm)
* `KeyManager` interface that can leveraged to manage/use keys (create, sign etc) as desired per the given use case. examples of concrete implementations include: AWS KMS, Azure Key Vault, Google Cloud KMS, Hashicorp Vault etc
* Concrete implementation of `KeyManager` that stores keys in memory

# Usage

## `dsa`

### Key Generation

the `dsa` package provides algorithm ID's that can be passed to the `GenerateKey` function i.e.

```go
package main

import (
  "fmt"        
  "github.com/abaxxtech/abaxx-id-go/internal/crypto/dsa"
)

func main() {
  privateJwk, err := dsa.GeneratePrivateKey(dsa.AlgorithmIDSECP256K1)
  if err != nil {
	  fmt.Printf("Failed to generate private key: %v\n", err)
	  return
  }
}
```

### Signing

Signing takes a private key and a payload to sign. e.g.

```go
package main

import (
  "fmt"
  "github.com/abaxxtech/abaxx-id-go/internal/crypto/dsa"
)

func main() {
  // Generate private key
  privateJwk, err := dsa.GeneratePrivateKey(dsa.AlgorithmIDSECP256K1)
  if err != nil {
    fmt.Printf("Failed to generate private key: %v\n", err)
    return
  }

  // Payload to be signed
  payload := []byte("hello world")

  // Signing the payload
  signature, err := dsa.Sign(payload, privateJwk)
  if err != nil {
    fmt.Printf("Failed to sign: %v\n", err)
    return
  }
}
```

### Verifying

Verifying takes a public key, the payload that was signed, and the signature. i.e.

```go
package main

import (
  "fmt"
  "github.com/abaxxtech/abaxx-id-go/internal/crypto/dsa"
)

func main() {
  // Generate ED25519 private key
  privateJwk, err := dsa.GeneratePrivateKey(dsa.AlgorithmIDED25519)
  if err != nil {
    fmt.Printf("Failed to generate private key: %v\n", err)
    return  
  }

  // Payload to be signed
  payload := []byte("hello world")

  // Sign the payload
  signature, err := dsa.Sign(payload, privateJwk)
  if err != nil {
    fmt.Printf("Failed to sign: %v\n", err)
    return
  }

  // Get the public key from the private key
  publicJwk := dsa.GetPublicKey(privateJwk)

  // Verify the signature
  legit, err := dsa.Verify(payload, signature, publicJwk)
  if err != nil {
    fmt.Printf("Failed to verify: %v\n", err)
    return
  }

  if !legit {
    fmt.Println("Failed to verify signature")
  } else {
    fmt.Println("Signature verified successfully")
  }
}
```

> [!NOTE]
> `ecdsa` and `eddsa` provide the same high level api as `dsa`, but specifically for algorithms within those respective families. this makes it so that if you add an additional algorithm, it automatically gets picked up by `dsa` as well.

# Directory Structure

```sh
crypto
├── README.md
├── doc.go
├── dsa
│   ├── README.md
│   ├── dsa.go
│   ├── dsa_test.go
│   ├── ecdsa
│   │   ├── ecdsa.go
│   │   ├── secp256k1.go
│   │   └── secp256k1_test.go
│   └── eddsa
│       ├── ed25519.go
│       └── eddsa.go
├── keymanager.go
└── keymanager_test.go
```
