# ID++ Protocol

<!-- @format -->

## Introduction

This repository contains ID++ a reference implementation of Decentralized Web Node (DWN) as per the [specification](https://identity.foundation/decentralized-web-node/spec/). This specification is in a draft state and very much so a WIP. For the foreseeable future, a lot of the work on DWN will be split across this repo and the repo that houses the specification.

Proposals and issues for the specification itself should be submitted as pull requests.

## Documentation

* [ID Dwn Whitepaper](docs/wp.pdf)

* [ID SDK Docs](https://dwnprotocol.gitbook.io/id++-sdk-docs)

## Usage

> `./bin/abaxx-id` is a CLI tool that allows you to create and manage DIDs, DWNs, and VCs.

## Integration Readme's

* [crypto](./internal/crypto/README.md)
* [dids](./internal/dids/README.md)
* [jwt](./internal/jwt/README.md)
* [jws](./internal/jws/README.md)
* [vc](./internal/vc/README.md)

## SQL Message Store

The SQL message store provides a persistent storage implementation for DWN messages. Here's how to use it:

### Installation

First, ensure you have a supported SQL database (e.g., PostgreSQL) installed and running.

### Basic Usage

```go
import (
    "github.com/abaxx/abaxx-id-go/pkg/store"
)

// Configure the message store
config := store.MessageStoreSQLConfig{
    DriverName:     "postgres",
    DataSourceName: "postgres://user:password@localhost:5432/dbname?sslmode=disable",
}

// Create a new message store instance
messageStore, err := store.NewMessageStoreSQL(config)
if err != nil {
    log.Fatal(err)
}

// Initialize the database schema
err = messageStore.Open()
if err != nil {
    log.Fatal(err)
}
defer messageStore.Close()

// Store a message
message := map[string]interface{}{
    "content": "Hello, World!",
    "type": "text/plain",
}
indexes := store.KeyValues{
    "interface": "MessagesInterface",
    "method": "Create",
    "dateCreated": time.Now().String(),
}
err = messageStore.Put("tenant1", message, indexes, nil)

// Query messages
filters := []store.Filter{
    {
        Property: "interface",
        Operator: "=",
        Value: "MessagesInterface",
    },
}
sort := &store.MessageSort{
    Property: "dateCreated",
    Direction: "DESC",
}
pagination := &store.Pagination{
    Limit: 10,
}
messages, cursor, err := messageStore.Query("tenant1", filters, sort, pagination, nil)
```

### Features

- CRUD operations for DWN messages
- Multi-tenant support
- Flexible querying with filters
- Sorting and pagination
- Support for message indexing
- CBOR encoding for efficient storage

### Database Schema

The message store automatically creates the necessary database schema when `Open()` is called. The schema includes:
- Message content storage
- Index fields for efficient querying
- Support for encoded data and metadata

For more details about the implementation, see [messagestore-sql.go](pkg/store/messagestore-sql.go).


