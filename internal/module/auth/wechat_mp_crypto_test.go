package auth

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/binary"
	"testing"
)

func pkcs7Pad(b []byte, blockSize int) []byte {
	pad := blockSize - (len(b) % blockSize)
	if pad == 0 {
		pad = blockSize
	}
	return append(b, bytes.Repeat([]byte{byte(pad)}, pad)...)
}

func encryptWeChatMPTextForTest(aesKey []byte, appID string, plain string) (string, error) {
	rand16 := bytes.Repeat([]byte{0x11}, 16)
	msg := []byte(plain)
	lenBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(lenBytes, uint32(len(msg)))
	payload := append(append(append(rand16, lenBytes...), msg...), []byte(appID)...)
	payload = pkcs7Pad(payload, aes.BlockSize)
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return "", err
	}
	iv := aesKey[:aes.BlockSize]
	out := make([]byte, len(payload))
	cipher.NewCBCEncrypter(block, iv).CryptBlocks(out, payload)
	return base64.StdEncoding.EncodeToString(out), nil
}

func TestDecryptWeChatMPTextWithKey_RoundTrip(t *testing.T) {
	aesKey := bytes.Repeat([]byte{0x22}, 32)
	appID := "wx1234567890abcd"
	plain := "<xml><FromUserName>openid</FromUserName><Event>SCAN</Event><EventKey>scene</EventKey></xml>"

	enc, err := encryptWeChatMPTextForTest(aesKey, appID, plain)
	if err != nil {
		t.Fatalf("encrypt failed: %v", err)
	}
	got, err := decryptWeChatMPTextWithKey(aesKey, appID, enc)
	if err != nil {
		t.Fatalf("decrypt failed: %v", err)
	}
	if got != plain {
		t.Fatalf("unexpected plaintext: %q != %q", got, plain)
	}
}

func TestDecryptWeChatMPTextWithKey_AppIDMismatch(t *testing.T) {
	aesKey := bytes.Repeat([]byte{0x33}, 32)
	enc, err := encryptWeChatMPTextForTest(aesKey, "wx_app_ok", "hello")
	if err != nil {
		t.Fatalf("encrypt failed: %v", err)
	}
	if _, err := decryptWeChatMPTextWithKey(aesKey, "wx_app_other", enc); err == nil {
		t.Fatalf("expected error on appid mismatch")
	}
}

