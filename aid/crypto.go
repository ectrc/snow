package aid

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
)

type keyPair struct {
	PrivateKey rsa.PrivateKey
	PublicKey  rsa.PublicKey
}

var KeyPair = GeneratePublicPrivateKeyPair()

func GeneratePublicPrivateKeyPair() keyPair {
	privateKey, _ := rsa.GenerateKey(rand.Reader, 2048)
	publicKey := privateKey.PublicKey

	return keyPair{
		PrivateKey: *privateKey,
		PublicKey:  publicKey,
	}
}

func (k *keyPair) EncryptAndSign(message []byte) ([]byte, []byte) {
	encryptedMessage, _ := rsa.EncryptPKCS1v15(rand.Reader, &k.PublicKey, message)
	signature, _ := rsa.SignPKCS1v15(rand.Reader, &k.PrivateKey, 0, encryptedMessage)

	return encryptedMessage, signature
}

func (k *keyPair) EncryptAndSignB64(message []byte) (string, string) {
	encryptedMessage, signature := k.EncryptAndSign(message)

	return Base64Encode(encryptedMessage), Base64Encode(signature)
}

func (k *keyPair) DecryptAndVerify(encryptedMessage []byte, signature []byte) []byte {
	decryptedMessage, _ := rsa.DecryptPKCS1v15(rand.Reader, &k.PrivateKey, encryptedMessage)
	_ = rsa.VerifyPKCS1v15(&k.PublicKey, 0, encryptedMessage, signature)

	return decryptedMessage
}

func (k *keyPair) ExportPrivateKey() []byte {
	privateKey := x509.MarshalPKCS1PrivateKey(&k.PrivateKey)
	return privateKey
}

func (k *keyPair) ExportPublicKey() []byte {
	publicKey := x509.MarshalPKCS1PublicKey(&k.PublicKey)
	return publicKey
}

func Base64Encode(input []byte) string {
	return base64.StdEncoding.EncodeToString(input)
}

func Base64Decode(input string) ([]byte, bool) {
	data, err := base64.StdEncoding.DecodeString(input)
	if err != nil {
			return []byte{}, true
	}
	return data, false
}

func Hash(input []byte) string {
	shaBytes := sha256.Sum256(input)
	return hex.EncodeToString(shaBytes[:])
}