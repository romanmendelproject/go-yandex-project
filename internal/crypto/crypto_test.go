package crypto

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCryptography(t *testing.T) {
	var metrics []byte
	key := "test"

	hash := GetHash(metrics, key)

	h := hmac.New(sha256.New, []byte(key))

	require.Equal(t, hash, base64.StdEncoding.EncodeToString(h.Sum(nil)))
}

func Example() {
	var metrics []byte
	key := "test"

	fmt.Println(GetHash(metrics, key))
	// Output:
	// rXEUjHnyGrnuxR6lx90rZoeS98DTU0rmayL3HGFSP7M=
}
