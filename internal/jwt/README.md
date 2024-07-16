# `jwt`

## Table of Contents
- [﻿Table of Contents](#table-of-contents) 
- [﻿Usage](#usage) 
    - [﻿Signing](#signing) 
    - [﻿Verifying](#verifying) 
- [﻿Directory Structure](#directory-structure) 
## Usage
### Signing
```go
package main

import (
 "fmt"
 "github.com/abaxxtech/abaxx-id-go/internal/dids/didjwk"
    "github.com/abaxxtech/abaxx-id-go/internal/jwt"
)

func main() { 
 did, err := didjwk.Create()
 if err != nil {
  panic(err)
 }

 claims := jwt.Claims{
  Issuer: did.URI,
  Misc:   map[string]interface{}{"c_nonce": "abcd123"},
 }

 jwt, err := jwt.Sign(claims, did)
 if err != nil {
  panic(err)
 }
}
```
## Verifying
```go
package main

import (
 "fmt"
 "github.com/abaxxtech/abaxx-id-go/internal/dids"
    "github.com/abaxxtech/abaxx-id-go/internal/jwt"
)

func main() {
    someJWT := "SOME_JWT"
 ok, err := jwt.Verify(signedJWT)
 if err != nil {
  panic(err)
 }

    if (!ok) {
        fmt.Printf("dookie JWT")
    }
}
```
specifying a specific category of key to use relative to the did provided can be done in the same way shown with `jws.Sign` 

### Directory Structure
```sh
jwt
├── jwt.go
└── jwt_test.go
```



<!-- eraser-additional-content -->
## Diagrams
<!-- eraser-additional-files -->
<a href="/internal/jwt/README-JWT Signing and Verifying Process-1.eraserdiagram" data-element-id="fYebsjTdni-srAwSgGMY1"><img src="/.eraser/7vYPAsAv9IcfykzVybNF___pHaokLkHewZxZhanJWMXDLMn78l2___---diagram----dd932827dd58043b071aa9dcf5ac809f-JWT-Signing-and-Verifying-Process.png" alt="" data-element-id="fYebsjTdni-srAwSgGMY1" /></a>
<!-- end-eraser-additional-files -->
<!-- end-eraser-additional-content -->
<!--- Eraser file: https://app.eraser.io/workspace/7vYPAsAv9IcfykzVybNF --->