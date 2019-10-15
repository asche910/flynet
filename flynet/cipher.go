package flynet

import (
	"crypto/aes"
	"crypto/cipher"
)

const key = "asche910-flynet-"

var (
	commonIV = []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f}
)

type Cipher struct {
	encoder cipher.Stream
	decoder cipher.Stream
	key     []byte
	method  string
}

func NewCipherInstance() *Cipher{
	c, err := aes.NewCipher([]byte(key))
	if err != nil {
		logger.Println("aes.NewCipher failed!", err)
	}
	enc := cipher.NewCFBEncrypter(c, commonIV)
	dec := cipher.NewCFBDecrypter(c, commonIV)
	return &Cipher{encoder:enc, decoder:dec}
}

func (cipher *Cipher) Encrypt(dst, src []byte) {
	cipher.encoder.XORKeyStream(dst, src)
}

func (cipher *Cipher) Decrypt(dst, src []byte) {
	cipher.decoder.XORKeyStream(dst, src)
}
