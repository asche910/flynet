package fly

import (
	"fmt"
	"testing"
)

const text = "hello,world!----++++++"
const text2 = "hello,world!worldworlds  sdf	 sfj  sa	 sfjslabb\n sfbzcvbxqeuqyuiyajhlfjdbz,v  sf"

func testCipher(t *testing.T, method string) {
	cipherInstance := NewCipherInstance(key, method)
	testSingleMethod(cipherInstance, text)
	testSingleMethod(cipherInstance, text2)
	testSingleMethod(cipherInstance, "safhj;		sfafajksfb")
	testSingleMethod(cipherInstance, "\ttextLen := len(text)\n\ttextLen := len(text)\n\ttext"+
		"Len := len(text)\n\ttextLen :\ttextLen := len(text)\n\ttext"+
		"Len := len(text)\n\ttextLen := len(text)\n\ttextLen := len(text)\n\ttextLen"+
		"Len := len(text)\n\ttextLen := len(text)\n\ttextLen := len(text)\n\ttextLen"+
		"Len := len(text)\n\ttextLen := len(text)\n\ttextLen := len(text)\n\ttextLen"+
		" := len(text)\n= len(text)\n\ttextLen := len(tex\ttextLen := len(text)\n"+
		"t)\n\ttextLen := len(text)\n")

	testSingleMethod(cipherInstance, "safhj;	sfsf	sfafajksfb")

}

func testSingleMethod(cipher *Cipher, raw string) bool {
	textLen := len(raw)
	decryptBuff := make([]byte, textLen)
	encryptBuff := make([]byte, textLen)
	cipher.Encrypt(encryptBuff, []byte(raw))
	cipher.Decrypt(decryptBuff, encryptBuff)
	if string(decryptBuff) != raw {
		fmt.Println("NOT PASS!!! ", string(decryptBuff))
		return false
	} else {
		fmt.Println("PASS")
		return true
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

func TestChacha20(t *testing.T) {
	testCipher(t, "chacha20")
	//testCipher(t, "chacha20-ietf")
}

func TestRC4MD5(t *testing.T) {
	testCipher(t, "rc4-md5")
	//testCipher(t, "rc4-md5-6")
}
