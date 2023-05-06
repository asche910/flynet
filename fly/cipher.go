package fly

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rc4"
	"fmt"

	"github.com/aead/chacha20"
)

var key = "asche910-flynet-"

type EncOrDec int

const (
	ENC EncOrDec = iota
	DEC
)

var (
	IV = []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f}
)

var CipherMap = map[string]*cipherEntity{
	"aes-128-cfb":   {16, 16, newAESCFBStream},
	"aes-192-cfb":   {24, 16, newAESCFBStream},
	"aes-256-cfb":   {32, 16, newAESCFBStream},
	"aes-128-ctr":   {16, 16, newAESCTRStream},
	"aes-192-ctr":   {24, 16, newAESCTRStream},
	"aes-256-ctr":   {32, 16, newAESCTRStream},
	"rc4-md5":       {16, 16, newRC4MD5Stream},
	"rc4-md5-6":     {16, 6, newRC4MD5Stream},
	"chacha20":      {32, 8, newChaCha20Stream},
	"chacha20-ietf": {32, 12, newChaCha20Stream},
}

// generate a fixed length key
func genKey(rawKey string, len int) []byte {
	key := make([]byte, 256)
	cur := 0
	for cur < len {
		sum := md5.Sum([]byte(fmt.Sprintf("%s%d", rawKey, cur)))
		copy(key[cur:cur+16], sum[:])
		cur += 16
	}
	return key[:len]
}

type Cipher struct {
	encoder cipher.Stream
	decoder cipher.Stream
	key     []byte
	method  string
}

func (cipher *Cipher) Encrypt(dst, src []byte) {
	cipher.encoder.XORKeyStream(dst, src)
}

func (cipher *Cipher) EncryptAndGet(src []byte) (dst []byte) {
	dst = make([]byte, len(src))
	cipher.encoder.XORKeyStream(dst, src)
	return
}

func (cipher *Cipher) Decrypt(dst, src []byte) {
	cipher.decoder.XORKeyStream(dst, src)
}

func (cipher *Cipher) DecryptAndGet(src []byte) (dst []byte) {
	dst = make([]byte, len(src))
	cipher.decoder.XORKeyStream(dst, src)
	return
}

// a detail encrypt or decrypt method
type cipherEntity struct {
	keyLen    int
	ivLen     int
	newStream func(key, iv []byte, eod EncOrDec) cipher.Stream
}

func NewCipherInstance(secretKey, method string) *Cipher {
	if secretKey == "" {
		secretKey = key
	}
	entity := CipherMap[method]
	if entity == nil {
		entity = CipherMap["aes-256-cfb"]
	}

	key := genKey(secretKey, entity.keyLen)
	newIV := IV[:entity.ivLen]
	enc := entity.newStream(key, newIV, ENC)
	dec := entity.newStream(key, newIV, DEC)
	return &Cipher{encoder: enc, decoder: dec}
}

func newAESCFBStream(key, iv []byte, eod EncOrDec) cipher.Stream {
	block, err := aes.NewCipher(key)
	if err != nil {
		logger.Println("aes.NewCipher failed --->", err)
		return nil
	}
	if eod == ENC {
		enc := cipher.NewCFBEncrypter(block, IV)
		return enc
	} else {
		dec := cipher.NewCFBDecrypter(block, IV)
		return dec
	}
}

func newAESCTRStream(key, iv []byte, eod EncOrDec) cipher.Stream {
	block, err := aes.NewCipher(key)
	if err != nil {
		logger.Println("aes.NewCipher failed --->", err)
		return nil
	}
	return cipher.NewCTR(block, iv)
}

func newChaCha20Stream(key, iv []byte, eod EncOrDec) cipher.Stream {
	stream, err := chacha20.NewCipher(iv, key)
	if err != nil {
		logger.Println("chacha20.NewCipher failed --->", err)
		return nil
	}
	return stream
}

func newRC4MD5Stream(key, iv []byte, eod EncOrDec) cipher.Stream {
	hs := md5.New()
	hs.Write(key)
	hs.Write(iv)
	rc4key := hs.Sum(nil)
	stream, err := rc4.NewCipher(rc4key)
	if err != nil {
		logger.Println("rc4.NewCipher failed --->", err)
		return nil
	}
	return stream
}
