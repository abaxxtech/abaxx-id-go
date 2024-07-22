package dids

import (
	"context"
	"sync"

	"github.com/abaxxtech/abaxx-id-go/internal/dids/did"
	"github.com/abaxxtech/abaxx-id-go/internal/dids/didcore"
	"github.com/abaxxtech/abaxx-id-go/internal/dids/diddht"
	"github.com/abaxxtech/abaxx-id-go/internal/dids/didion"
	"github.com/abaxxtech/abaxx-id-go/internal/dids/didjwk"
	"github.com/abaxxtech/abaxx-id-go/internal/dids/didweb"
)

// Resolve resolves the provided DID URI. This function is capable of resolving
func Resolve(uri string) (didcore.ResolutionResult, error) {
	return getDefaultResolver().Resolve(uri)
}

// ResolveWithContext resolves the provided DID URI. This function is capable of resolving
func ResolveWithContext(ctx context.Context, uri string) (didcore.ResolutionResult, error) {
	return getDefaultResolver().ResolveWithContext(ctx, uri)
}

var (
	instance *didResolver
	once     sync.Once
)

func getDefaultResolver() *didResolver {
	once.Do(func() {
		instance = &didResolver{
			resolvers: map[string]didcore.MethodResolver{
				"dht": diddht.DefaultResolver(),
				"jwk": didjwk.Resolver{},
				"web": didweb.Resolver{},
				"ion": didion.Resolver{},
			},
		}
	})

	return instance
}

type didResolver struct {
	resolvers map[string]didcore.MethodResolver
}

func (r *didResolver) Resolve(uri string) (didcore.ResolutionResult, error) {
	return r.resolveWithContext(context.Background(), uri)
}

func (r *didResolver) ResolveWithContext(ctx context.Context, uri string) (didcore.ResolutionResult, error) {
	return r.resolveWithContext(ctx, uri)
}

func (r *didResolver) resolveWithContext(ctx context.Context, uri string) (didcore.ResolutionResult, error) {
	did, err := did.Parse(uri)
	if err != nil {
		return didcore.ResolutionResultWithError("invalidDid"), didcore.ResolutionError{Code: "invalidDid"}
	}

	resolver, found := r.resolvers[did.Method]
	if !found {
		return didcore.ResolutionResultWithError("methodNotSupported"), didcore.ResolutionError{Code: "methodNotSupported"}
	}

	return resolver.ResolveWithContext(ctx, uri)
}
