package ss58_test

import (
	"github.com/itering/subscan/util"
	"github.com/itering/subscan/util/base58"
	"github.com/itering/subscan/util/ss58"
	"github.com/prometheus/common/log"
	"testing"
)

func TestDecode(t *testing.T) {
	address := "5FcEGUiujfdWyf6RME1G8pCTkmkgXFDECaTSpVDWVnNiZJXR"
	decoded := ss58.Decode(address, 42)
	log.Warn(base58.Encode(util.HexToBytes(decoded)))
	if decoded != "9cbfadc7579a27fcb3ea4bb1940aade652d1dd9a2dc69c9920f1de42d8ca0234" {
		t.Fail()
	}
}

func TestEncode(t *testing.T) {
	address := "0x88b3bfe1410ed8a12cd8a2c230e97cfd5a9fb1cc95ac859ec9c9a2ecfe7cf84f"
	encoded := ss58.Encode(address, 2)
	if encoded != "FfZRiEyrJwgxFZx1QsCnDjaJCHXoeUS4v4Hs1Yo8GpVveNQ" {
		t.Fail()
	}
}

func TestKusamaDecode(t *testing.T) {
	//address := "GTug9rrdeBDadKXJ9DM5pUgYKnWJeVrvY8UeWpm3PpQgRq9"
	address := "5FGDTP2nHiUzcQ4zCto73TvBKC7GerbyXTjAHp9Yruypoq6U"
	address1 := "14E5nqKAp3oAJcmzgZhUD2RcptBeUBScxKHgJKU4HPNcKVf3"

	ss58Format := base58.Decode(address)
	log.Warn(util.BytesToHex(ss58Format))
	decoded := ss58.Decode(address, 42)
	log.Warn("Decoded: " + decoded)

	log.Warn(util.BytesToHex(base58.Decode(address1)))
	log.Warn(ss58.Decode(address1, 0))
}
