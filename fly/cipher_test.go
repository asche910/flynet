package fly

import (
	"fmt"
	"testing"
)

const text = "hello,world!"

func testCipher(t *testing.T, method string) {
	textLen := len(text)
	cipherInstance := NewCipherInstance(key, method)

	textBuff := make([]byte, textLen)
	encryptBuff := make([]byte, textLen)

	cipherInstance.Encrypt(encryptBuff, []byte(text))
	fmt.Println("after", encryptBuff)
	cipherInstance.Decrypt(textBuff, encryptBuff)
	if string(textBuff) != text {
		t.Error(method, " test failed!")
	}
}

func TestCipher(t *testing.T) {
	key := "123456"
	keyLen := len(key)
	buff := make([]byte, 1024)
	cipherInstance := NewCipherInstance(key, "aes-128-cfb")
	cipherInstance.encoder.XORKeyStream(buff, []byte(key))
	fmt.Println("before", string(buff[:keyLen]))
	cipherInstance.decoder.XORKeyStream(buff, buff[:keyLen])
	fmt.Println("after", string(buff[:keyLen]))
}

func TestAESCTR(t *testing.T) {
	testCipher(t, "aes-128-ctr")
}

func TestChacha20(t *testing.T)  {
	testCipher(t, "chacha20")
	//testCipher(t, "chacha20-ietf")
}

func TestRC4MD5(t *testing.T) {
	testCipher(t, "rc4-md5")
	//testCipher(t, "rc4-md5-6")
}