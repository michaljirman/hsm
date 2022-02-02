package main

import (
	"crypto/rand"
	"log"

	"github.com/ThalesIgnite/crypto11"
)

func randomBytes() []byte {
	result := make([]byte, 32)
	rand.Read(result)
	return result
}

const rsaSize = 2048

func main() {
	ctx, err := crypto11.Configure(&crypto11.Config{
		Path:              "/opt/nfast/toolkits/pkcs11/libcknfast.so",
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
