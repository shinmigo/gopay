package kernel

import (
	"bytes"
	"crypto"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"io/ioutil"
	"os"
	"strings"
)

/**
 *	解析PKCS8格式私钥
 */
func ParsePrivateKey(merchantPrivateKeyPath string) (*rsa.PrivateKey, error) {
	if len(merchantPrivateKeyPath) == 0 {
		return nil, errors.New(MerchantPrivateKeyEmpty)
	}
	byteContent, err := getCertContent(merchantPrivateKeyPath)
	if err != nil {
		return nil, err
	}
	byteContentLen := len(byteContent)
	if byteContentLen == 0 {
		return nil, errors.New(CertEmpty)
	}

	data := formatPKCS8PrivateKey(string(byteContent))
	if data == nil || len(data) == 0 {
		return nil, errors.New(PrivateKeyWrongFormat)
	}
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, errors.New(PrivateKeyWrongFormat)
	}
	privateKeyInter, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	privateKey, ok := privateKeyInter.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New(PrivateKeyWrongFormat)
	}

	return privateKey, nil
}

/**
 *	解析支付宝公钥
 */
func ParsePublicKey(aliPayPublicKeyPath string) (*rsa.PublicKey, error) {
	if len(aliPayPublicKeyPath) == 0 {
		return nil, errors.New(MerchantPrivateKeyEmpty)
	}
	byteContent, err := getCertContent(aliPayPublicKeyPath)
	if err != nil {
		return nil, err
	}
	byteContentLen := len(byteContent)
	if byteContentLen == 0 {
		return nil, errors.New(CertEmpty)
	}

	//格式化支付宝公钥
	data := formatPrivatePublicKey(string(byteContent), AliPayPublicKeyPrefix, AliPayPublicKeySuffix, 64)
	if data == nil || len(data) == 0 {
		return nil, errors.New(PublicKeyWrongFormat)
	}
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, errors.New(PublicKeyWrongFormat)
	}
	publicKeyInter, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	publicKey, ok := publicKeyInter.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New(PublicKeyWrongFormat)
	}

	return publicKey, nil
}

/**
 * 提取商户支付宝应用公钥证书序列号
 */
func GetMerchantCertSN(certPath string) (string, error) {
	byteContent, err := getCertContent(certPath)
	if err != nil {
		return "", err
	}
	byteContentLen := len(byteContent)
	if byteContentLen == 0 {
		return "", errors.New(CertEmpty)
	}

	cert, err := ParseAliPayCert(byteContent)
	if err != nil {
		return "", err
	}
	return getCertSN(cert), nil
}

/**
 * 提取支付宝公钥证书序列号
 */
func GetAliPayCertSN(certPath string) (string, *rsa.PublicKey, error) {
	byteContent, err := getCertContent(certPath)
	if err != nil {
		return "", nil, err
	}
	byteContentLen := len(byteContent)
	if byteContentLen == 0 {
		return "", nil, errors.New(CertEmpty)
	}

	cert, err := ParseAliPayCert(byteContent)
	if err != nil {
		return "", nil, err
	}
	publicKey, ok := cert.PublicKey.(*rsa.PublicKey)
	if ok == false {
		return "", nil, err
	}

	return getCertSN(cert), publicKey, nil
}

/**
 * 提取支付宝根证书序列号
 */
func GetAliPayRootCertSN(certPath string) (string, error) {
	byteContent, err := getCertContent(certPath)
	if err != nil {
		return "", err
	}
	byteContentLen := len(byteContent)
	if byteContentLen == 0 {
		return "", errors.New(CertEmpty)
	}

	certSnSlice := make([]string, 0, byteContentLen)
	aliPayRootCertSlice := strings.Split(string(byteContent), AliPayRootCertEnd)
	for _, certContent := range aliPayRootCertSlice {
		certContent = certContent + AliPayRootCertEnd

		cert, err := ParseAliPayCert([]byte(certContent))
		if err != nil {
			return "", err
		}
		if cert.SignatureAlgorithm == x509.SHA256WithRSA || cert.SignatureAlgorithm == x509.SHA1WithRSA {
			certSnSlice = append(certSnSlice, getCertSN(cert))
		}
	}
	aliPayRootCertCn := strings.Join(certSnSlice, "_")

	return aliPayRootCertCn, nil
}

/**
 *	获取证书文件内容
 */
func getCertContent(certPath string) ([]byte, error) {
	fileHandler, err := os.Open(certPath)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = fileHandler.Close()
	}()

	byteContent, err := ioutil.ReadAll(fileHandler)
	if err != nil {
		return nil, err
	}

	return byteContent, nil
}

/**
 *	获取序列号
 */
func getCertSN(cert *x509.Certificate) string {
	value := md5.Sum([]byte(cert.Issuer.String() + cert.SerialNumber.String()))
	return hex.EncodeToString(value[:])
}

/**
 *	解析支付宝证书
 */
func ParseAliPayCert(certContentByte []byte) (*x509.Certificate, error) {
	block, _ := pem.Decode(certContentByte)
	if block == nil {
		return nil, errors.New(CertWrongFormat)
	}

	return x509.ParseCertificate(block.Bytes)
}

/**
 *	格式化 PKCS1 私钥
 */
func formatPKCS1PrivateKey(privateKey string) []byte {
	privateKey = strings.Replace(privateKey, AliPayPKCS8Prefix, "", 1)
	privateKey = strings.Replace(privateKey, AliPayPKCS8Suffix, "", 1)

	return formatPrivatePublicKey(privateKey, AliPayPKCS1Prefix, AliPayPKCS1Suffix, 64)
}

/**
 *	格式化 PKCS8 私钥
 */
func formatPKCS8PrivateKey(privateKey string) []byte {
	privateKey = strings.Replace(privateKey, AliPayPKCS1Prefix, "", 1)
	privateKey = strings.Replace(privateKey, AliPayPKCS1Suffix, "", 1)

	return formatPrivatePublicKey(privateKey, AliPayPKCS8Prefix, AliPayPKCS8Suffix, 64)
}

/**
 *	格式化公钥|私钥
 */
func formatPrivatePublicKey(key, prefix, suffix string, lineCount int) []byte {
	if len(key) == 0 {
		return nil
	}

	key = strings.Replace(key, prefix, "", 1)
	key = strings.Replace(key, suffix, "", 1)
	var pendingString = []string{" ", "\n", "\r", "\t"}
	for i := range pendingString {
		value := pendingString[i]
		key = strings.Replace(key, value, "", -1)
	}
	formatKeyLen := len(key)
	count := formatKeyLen / lineCount
	if formatKeyLen%lineCount > 0 {
		count = count + 1
	}

	var buffers bytes.Buffer
	buffers.WriteString(prefix + "\n")
	for i := 0; i < count; i++ {
		a := i * lineCount
		b := a + lineCount
		if b > formatKeyLen {
			buffers.WriteString(key[a:])
		} else {
			buffers.WriteString(key[a:b])
		}
		buffers.WriteString("\n")
	}
	buffers.WriteString(suffix)

	return buffers.Bytes()
}

/**
 *	使用sha256签名
 */
func Sign(content []byte, merchantPrivateKey *rsa.PrivateKey, signType string) ([]byte, error) {
	if signType == AliPaySignType {
		hashHandler := sha1.New()
		hashHandler.Write(content)
		return rsa.SignPKCS1v15(rand.Reader, merchantPrivateKey, crypto.SHA1, hashHandler.Sum(nil))
	} else {
		hashHandler := sha256.New()
		hashHandler.Write(content)
		return rsa.SignPKCS1v15(rand.Reader, merchantPrivateKey, crypto.SHA256, hashHandler.Sum(nil))
	}
}

/**
 *	支付宝签名验证
 */
func Verify(src, sign []byte, key *rsa.PublicKey, signType string) error {
	if signType == AliPaySignType {
		hashHandler := sha1.New()
		hashHandler.Write(src)
		return rsa.VerifyPKCS1v15(key, crypto.SHA1, hashHandler.Sum(nil), sign)
	} else {
		hashHandler := sha256.New()
		hashHandler.Write(src)
		return rsa.VerifyPKCS1v15(key, crypto.SHA256, hashHandler.Sum(nil), sign)
	}
}
