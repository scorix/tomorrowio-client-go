package tomorrowio

import "github.com/google/wire"

var Set = wire.NewSet(
	NewClient,
)
