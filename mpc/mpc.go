package main

import (
	"crypto/rand"
	`errors`
	"log"
	`os`

	"github.com/ThalesIgnite/crypto11"
)

func randomBytes() []byte {
	result := make([]byte, 32)
	rand.Read(result)
	return result
}

const (
	rsaSize = 2048
	path = "/nfast/toolkits/pkcs11/libcknfast.so"
)

func main() {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		log.Fatal("file does not exists", err)
	}
	ctx, err := crypto11.Configure(&crypto11.Config{
		Path:              path,
		TokenSerial:       "6D30-03E0-D947",
		LoginNotSupported: true,
	})
	if err != nil {
		log.Fatal(err)
	}

	id := randomBytes()
	_, err = ctx.GenerateRSAKeyPair(id, rsaSize)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("RSA key pair was successfully generated")
}
