package rsa_tools

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
)

type Type int64

const (
	PKCS1 Type = iota
	PKCS8
)

type pkcsClient struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

type Cipher interface {
	Sign(src []byte, hash crypto.Hash) ([]byte, error)
	Verify(src []byte, sign []byte, hash crypto.Hash) error
}

func (this *pkcsClient) Sign(src []byte, hash crypto.Hash) ([]byte, error) {
	h := hash.New()
	h.Write(src)
	hashed := h.Sum(nil)
	return rsa.SignPKCS1v15(rand.Reader, this.privateKey, hash, hashed)
}

func (this *pkcsClient) Verify(src []byte, sign []byte, hash crypto.Hash) error {
	h := hash.New()
	h.Write(src)
	hashed := h.Sum(nil)
	return rsa.VerifyPKCS1v15(this.publicKey, hash, hashed, sign)
}

func New(privateKey []byte, publicKey []byte, privateKeyType Type) (Cipher, error) {

	blockPri, _ := pem.Decode([]byte(privateKey))
	if blockPri == nil {
		return nil, errors.New("private key error")
	}

	blockPub, _ := pem.Decode([]byte(publicKey))
	if blockPub == nil {
		return nil, errors.New("public key error")
	}

	priKey, err := genPriKey(blockPri.Bytes, privateKeyType)
	if err != nil {
		return nil, err
	}
	pubKey, err := genPubKey(blockPub.Bytes)
	if err != nil {
		return nil, err
	}
	return &pkcsClient{privateKey: priKey, publicKey: pubKey}, nil
}

func genPubKey(publicKey []byte) (*rsa.PublicKey, error) {
	pub, err := x509.ParsePKIXPublicKey(publicKey)
	if err != nil {
		return nil, err
	}
	return pub.(*rsa.PublicKey), nil
}

func genPriKey(privateKey []byte, privateKeyType Type) (*rsa.PrivateKey, error) {
	var priKey *rsa.PrivateKey
	var err error
	switch privateKeyType {
	case PKCS1:
		{
			priKey, err = x509.ParsePKCS1PrivateKey([]byte(privateKey))
			if err != nil {
				return nil, err
			}
		}
	case PKCS8:
		{
			prkI, err := x509.ParsePKCS8PrivateKey([]byte(privateKey))
			if err != nil {
				return nil, err
			}
			priKey = prkI.(*rsa.PrivateKey)
		}
	default:
		{
			return nil, errors.New("unsupport private key type")
		}
	}
	return priKey, nil
}
