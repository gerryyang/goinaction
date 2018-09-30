package main

import (
	"crypto"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"rsa_tools"
)

type ENCODE_MODE int64

const (
	HEX ENCODE_MODE = iota
	BASE64
)

var cipher rsa_tools.Cipher

// PKCS8
var privateKey = []byte(`
-----BEGIN PRIVATE KEY-----
MIIEvwIBADANBgkqhkiG9w0BAQEFAASCBKkwggSlAgEAAoIBAQDh41jjFnxId/ov
bYgTOp1BQd2EEn1pGOckU69eEZltnAvJ9ohfJvgqFagVQYAfuOJMJH5cd3xWpUq0
3YaDN0uqe4meBLz1pq8VF5pd+Wj0w3vkiaW7DvpsJGTop1kTwFCanMbLVXBoyukH
BPRdnFXGMit28xoySJFYHmv2/IqBs+ys4QCg1fWzPGLzem78ijzngWBvkNGmK+UD
A9YpCxGusvhZCs9ftOyA3aD3ojFcM8t6kyUOW+GLygD5d0eJltq4Z9WoS5XzzTJZ
VMc40gwVuSpdNp8mZxy6a84H3S9MPYg4HCRM9FVEZQzRJSVTGXKG4yhJLdxuvV5m
f0eDFb0TAgMBAAECggEBAM8H8hXgK/S3keQaPZdyJ2MCHSbJU4wZuO/Ai4BqHPcr
CFsIy6B6NQVNaApjSCzK5Q3ofK//C0TWpgvy5TAqY/1S0KS1rwJuzRVF1sO+rgV3
jXu+9NjnN3oaOSpLBwdlQfsTKdh+7FH/d2hpkBakDLFklhWlZiMkA6KB724ltER+
jFfBl1O210w9AM2aOmYKPuz82iIfcabWO/1T+d99cim6MODaHC/RcAadvWvzMrDU
q+6pn+dmuIe8vSD8DN6fJqHsGow59B1af5tcAgVVFAURzli0BqyLJSWp7roRomcg
S58eDBCzJSlOL4rBW5cbS+JPSkNu/RjvMzJQk5tR6EECgYEA/NRnWu0eyZv/dXSK
ktElyMB6uVjw3tpQN34sNZLXWoNL3pu3P539r9l7Y4JehSaKS7U2QL19O/HkyRDg
w+HY+WdEKqD+6PYEEIXn19UZvZbWV0pl7CwVWRYIQQKppWU+zhq4cFmI7wAidNGK
kCE8Wt+UAeZRsNAjv2IokkLTjhECgYEA5Lh1oLxKNvZ2TEv1lGEtAoF7sLu8U1zv
tKE3/Eor1RM6L2uI6sacMOqHKDwXbj2GHxYnobscEiBQEi7E6SoHwq1h6HNjwdXb
Wi2pakydhvDPOM4ETBu3jwMEac7O9Dfo9CGglaLwj7Md8hWKaU+8gEc9OurUMJBS
OVQFiOJehOMCgYEA7cy13a7TW1svnoDb6ZVwDW8Evxopi+IYuukgmc8gYNDHZnxd
kid+uYw74u93CZOjVev+OExB40T0JC2MypC9LG91jQbaW7ExR307ACU+TbT2qymd
zdH0zlLLtqHTgG5G8UHuojWEdw9QWUHRKxknlG1f352Kzlwmk1a2xEK4ipECgYBI
Q0dWy2afSutBW9ZxVOqFmidcRVRQ+lH5vd4UZdLHdVWy2cTeeHWstsyRF7tHZ0TS
2YsX/Cf4SiFCPWiVSmQ9S85dROfFvC2bpkWagi5bDgZKqjyNV0x9cLSaQW79lhSR
3XYBEQP0QuE5NTkP4NNrrBZaYQs9dLulxTgicXLvhQKBgQDRAcdE0hedSD7yfbjs
U0ZuMRWwoCSB0ZxIp0M67d2S3Mcdsw1/lw4lAtzJ3OHesfGpIFNrCC/ZI8x88KVL
uN9nDdPtCWJLx/4YYQ8JheOJOkRazYyZJom4l4Kmc6/TfDWrlGbMe3cr/fiBUDfM
PpMCLbgCqT07n+dWz7TShah3ZQ==
-----END PRIVATE KEY-----
`)

var publicKey = []byte(`
-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA4eNY4xZ8SHf6L22IEzqdQUHdhBJ9aRjnJFOvXhGZbZwLyfaIXyb4KhWoFUGAH7jiTCR+XHd8VqVKtN2GgzdLqnuJngS89aavFReaXflo9MN75Imluw76bCRk6KdZE8BQmpzGy1VwaMrpBwT0XZxVxjIrdvMaMkiRWB5r9vyKgbPsrOEAoNX1szxi83pu/Io854Fgb5DRpivlAwPWKQsRrrL4WQrPX7TsgN2g96IxXDPLepMlDlvhi8oA+XdHiZbauGfVqEuV880yWVTHONIMFbkqXTafJmccumvOB90vTD2IOBwkTPRVRGUM0SUlUxlyhuMoSS3cbr1eZn9HgxW9EwIDAQAB
-----END PUBLIC KEY-----
`)

func main() {
	cipher, err := rsa_tools.New(privateKey, publicKey, rsa_tools.PKCS8)
	if err != nil {
		fmt.Println(err)
	}

	plain_text := "action=query&channel=wechat&ts=1538231718&user_id=oEIpN5c8e34o6jaV5KG48vJTDpBA&fTcCu9Qi5QtjaarYkXS4u1zTxtF4igXX"

	signBytes, err := cipher.Sign([]byte(plain_text), crypto.SHA256)
	if err != nil {
		fmt.Println(err)
	}

	var enc ENCODE_MODE = BASE64

	if enc == BASE64 {
		sign_enc := base64.StdEncoding.EncodeToString(signBytes)
		fmt.Println(sign_enc)
		sign_dec, err := base64.StdEncoding.DecodeString(sign_enc)

		err = cipher.Verify([]byte(plain_text), sign_dec, crypto.SHA256)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("verify success")

	} else if enc == HEX {
		sign_enc := hex.EncodeToString(signBytes)
		fmt.Println(sign_enc)
		sign_dec, err := hex.DecodeString(sign_enc)

		err = cipher.Verify([]byte(plain_text), sign_dec, crypto.SHA256)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("verify success")
	}
}
