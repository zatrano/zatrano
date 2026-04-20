// Package graphql mounts a schema-first gqlgen server on Fiber (see Register) and provides
// per-request DataLoader attachment via WithLoaders / LoadersFrom.
//
//go:generate go run github.com/99designs/gqlgen@v0.17.78 generate --config ../../gqlgen.yml
package graphql
