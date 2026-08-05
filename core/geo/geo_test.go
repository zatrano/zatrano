package geo_test

import (
	"testing"

	"github.com/zatrano/framework/core/geo"
)

func TestGeoLookup(t *testing.T) {
	r := geo.New()
	loc := r.Lookup("8.8.8.8")
	if loc.CountryCode != "US" {
		t.Fatalf("%+v", loc)
	}
	loop := r.Lookup("127.0.0.1")
	if loop.CountryCode != "LO" {
		t.Fatalf("%+v", loop)
	}
	priv := r.Lookup("10.0.0.5")
	if priv.CountryCode != "PR" {
		t.Fatalf("%+v", priv)
	}
}
