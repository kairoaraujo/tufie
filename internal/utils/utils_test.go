package utils

import (
	"fmt"
	"testing"

	"github.com/kairoaraujo/tufie/internal/tuf"

	"github.com/stretchr/testify/assert"
)

func TestStringSha(t *testing.T) {

	str := StringSha("test string")
	assert.Equal(t, str, "d5579c46dfcc7f18207013e65b44e4cb4e2c2298f4ac457ba8f82743f31e930b")

}

func TestEncodeTrustedRoot(t *testing.T) {
	stringBase64 := EncodeTrustedRoot([]byte("123"))

	assert.Equal(t, "MTIz", stringBase64)
}

func TestDecodeTrustedRoot(t *testing.T) {

	rootMd, err := tuf.LoadTrustedRoot("../../tests/test-root.json")
	if err != nil {
		fmt.Println(err.Error())
	}
	rootB, _ := rootMd.ToBytes(false)
	stringBase64 := EncodeTrustedRoot(rootB)

	decodedRootMd := DecodeTrustedRoot(stringBase64)

	assert.Equal(t, rootMd, decodedRootMd)

}
