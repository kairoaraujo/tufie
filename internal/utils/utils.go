package utils

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/rdimitrov/go-tuf-metadata/metadata"
	"github.com/spf13/cobra"
)

func StringSha(s string) string {
	h := sha256.New()
	h.Write([]byte(s))
	bs := h.Sum(nil)

	return fmt.Sprintf("%x", bs)
}

func DecodeTrustedRoot(base64Root string) *metadata.Metadata[metadata.RootType] {
	bytes, err := base64.StdEncoding.DecodeString(base64Root)
	cobra.CheckErr(err)

	var rootJSON map[string]interface{}
	err = json.Unmarshal(bytes, &rootJSON)
	cobra.CheckErr(err)

	rootBytes, _ := json.MarshalIndent(rootJSON, "", " ")
	root, err := metadata.Root().FromBytes(rootBytes)
	cobra.CheckErr(err)

	return root

}

func EncodeTrustedRoot(root []byte) string {

	rootBase64 := base64.StdEncoding.EncodeToString(root)
	return rootBase64
}
