package service

import (
	"github.com/itering/substrate-api-rpc"
	"testing"
)

func TestDecodeEnv(t *testing.T) {
	event := ""
	metadataInstant := s.getMetadataInstant(spec)
	substrate.DecodeEvent(event, metadataInstant, spec)
}
