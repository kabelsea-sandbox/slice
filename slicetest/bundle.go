package slicetest

import (
	"github.com/kabelsea-sanbox/slice"
)

// NewTestBundle
func NewTestBundle(bundle slice.Bundle) *TestBundle {
	return &TestBundle{
		bundle: bundle,
	}
}

type TestBundle struct {
	bundle slice.Bundle
}
