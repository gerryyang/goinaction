package common_utils

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"sort"
)

// 入参：
// appkey: 加盐参数
// privateKey: 计算签名的私钥
// params: 
func CalcBillingSign(appKey string, privateKey string, params map[string]interface{}) string {
	if params == nil {
		return ""
	}

	var sortedParams string

	// 1. 组合请求参数，字典序升序排列拼接成字符串
	keys := make([]string, 0)
	for k, _ := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		sortedParams += fmt.Sprintf("%s=%v&", k, params[k])
	}
	sortedParams = sortedParams[:(len(sortedParams) - 1)]

	// 2. 拼接appkey
	sortedParams += appKey

	// 3. 对 sortedParams 进行RSA-SHA256签名
	// 3.1 解码私钥
	prikeyBytes, err := base64.StdEncoding.DecodeString(privateKey)
	if err != nil {
		fmt.Println("base64.StdEncoding.DecodeString failed! error info is:[%+v]", err)
		return ""
	}
	// 3.2 读取私钥
	priKeyInf, err := x509.ParsePKCS8PrivateKey(prikeyBytes)
	if err != nil {
		fmt.Println("x509.ParsePKCS8PrivateKey failed! error info is:[%+v]", err)
		return ""
	}
	priKey := priKeyInf.(*rsa.PrivateKey)

	// 4. sha256 摘要params
	shaBytes := sha256.Sum256([]byte(sortedParams))

	// 5. 私钥签名
	signature, err := rsa.SignPKCS1v15(rand.Reader, priKey, crypto.SHA256, shaBytes[:])
	if err != nil {
		fmt.Println("rsa.SignPKCS1v15 failed! error info is:[%+v]", err)
		return ""
	}

	// 6. base64encode
	return base64.StdEncoding.EncodeToString(signature)
}
 
