<p><a target="_blank" href="https://app.eraser.io/workspace/VplWee70eSRvVNyZOOrd" id="edit-in-eraser-github-link"><img alt="Edit in Eraser" src="https://firebasestorage.googleapis.com/v0/b/second-petal-295822.appspot.com/o/images%2Fgithub%2FOpen%20in%20Eraser.svg?alt=media&amp;token=968381c8-a7e7-472a-8ed6-4a6626da5501"></a></p>

# `jws` 
## Table of Contents 
- [﻿Features](#features) 
- [﻿Usage](#usage) 
    - [﻿Signing:](#signing) 
    - [﻿Detached Content](#detached-content) 
    - [﻿Verifying](#verifying) 
    - [﻿Directory Structure](#directory-structure) 
### Features
- Signing a JWS (JSON Web Signature) with a DID
- Verifying a JWS with a DID
## Usage
### Signing
```go
package main

import (
    "fmt"
    "github.com/abaxxtech/abaxx-id-go/internal/dids/didjwk"
    "github.com/abaxxtech/abaxx-id-go/internal/jws"
)

func main() {
    did, err := didjwk.Create()
    if err != nil {
        fmt.Printf("failed to create did: %v", err)
        return
    }

    payload := map[string]interface{}{"hello": "world"}
    
    compactJWS, err := jws.Sign(payload, did)
    if err != nil {
        fmt.Printf("failed to sign: %v", err)
        return
    }

    fmt.Printf("compact JWS: %s", compactJWS)
}
```
## Detached Content
returning a JWS with detached content can be done like so:

```go
package main

import (
    "fmt"
    "github.com/abaxxtech/abaxx-id-go/internal/dids/didjwk"
    "github.com/abaxxtech/abaxx-id-go/internal/jws"
)

func main() {
    did, err := didjwk.Create()
    if err != nil {
        fmt.Printf("failed to create did: %v", err)
        return
    }

    payload := map[string]interface{}{"hello": "world"}
    
    compactJWS, err := jws.Sign(payload, did, Detached(true))
    if err != nil {
        fmt.Printf("failed to sign: %v", err)
        return
    }

    fmt.Printf("compact JWS: %s", compactJWS)
}
```
specifying a specific category of key associated with the provided did to sign with can be done like so:

```go
package main

import (
    "fmt"
    "github.com/abaxxtech/abaxx-id-go/internal/dids/didjwk"
    "github.com/abaxxtech/abaxx-id-go/internal/jws"
)

func main() {
    bearerDID, err := didjwk.Create()
    if err != nil {
        fmt.Printf("failed to create did: %v", err)
        return
    }

    payload := map[string]interface{}{"hello": "world"}
    
    compactJWS, err := jws.Sign(payload, did, Purpose("authentication"))
    if err != nil {
        fmt.Printf("failed to sign: %v", err)
    }

    fmt.Printf("compact JWS: %s", compactJWS)
}
```
### Verifying
```go
package main

import (
    "fmt"
    "github.com/abaxxtech/abaxx-id-go/internal/dids/didjwk"
    "github.com/abaxxtech/abaxx-id-go/internal/jws"
)

func main() {
    compactJWS := "SOME_JWS"
    ok, err := jws.Verify(compactJWS)
    if (err != nil) {
        fmt.Printf("failed to verify JWS: %v", err)
    }

    if (!ok) {
        fmt.Errorf("integrity check failed")
    }
}
```
>  an error is returned if something in the process of verification failed whereas `!ok` means the signature is actually shot 

### Directory Structure
```sh
jws
├── jws.go
└── jws_test.go
```



<!-- eraser-additional-content -->
## Diagrams
<!-- eraser-additional-files -->
<a href="/internal/jws/README-JWS Signing and Verification Process-1.eraserdiagram" data-element-id="146dvQmbwYPBpGUM2RzBk"><img src="/.eraser/VplWee70eSRvVNyZOOrd___pHaokLkHewZxZhanJWMXDLMn78l2___---diagram----cd2c3da8ebf4159ba46b99cd1225e9e6-JWS-Signing-and-Verification-Process.png" alt="" data-element-id="146dvQmbwYPBpGUM2RzBk" /></a>
<!-- end-eraser-additional-files -->
<!-- end-eraser-additional-content -->
<!--- Eraser file: https://app.eraser.io/workspace/VplWee70eSRvVNyZOOrd --->