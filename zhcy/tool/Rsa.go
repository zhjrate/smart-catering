package tool

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"log"
	"os"
)

type ToolRsa struct {
	PrivateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
	Code       bool
}

func (r *ToolRsa) CreateRSAKeyPair() error {
	// 生成RSA密钥对
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}
	publicKey := &privateKey.PublicKey
	r.PublicKey = publicKey
	r.PrivateKey = privateKey
	return nil
}
func (r ToolRsa) VerifySignData(signatureBase64 string, originalData []byte) bool {
	// 解码 Base64 编码的签名
	signature, err := base64.StdEncoding.DecodeString(signatureBase64)
	if err != nil {
		log.Println("解码签名失败:", err)
		return false
	}
	// 对原始数据进行哈希
	hasher := sha256.New()
	hasher.Write(originalData)
	hashed := hasher.Sum(nil)

	// 使用公钥验证签名
	err = rsa.VerifyPKCS1v15(r.PublicKey, crypto.SHA256, hashed, signature)
	if err != nil {
		log.Println("验证签名失败:", err)
		return false
	}
	// 签名验证成功
	return true
}

func (r ToolRsa) SignRSAdata(data []byte) (string, error) {
	hasher := sha256.New()
	hasher.Write(data)
	hashed := hasher.Sum(nil)

	// 使用RSA私钥对消息摘要进行签名
	signature, err := rsa.SignPKCS1v15(rand.Reader, r.PrivateKey, crypto.SHA256, hashed)
	if err != nil {
		log.Fatal("Failed to sign data:", err)
		return "", err
	}

	// 将签名转换为Base64编码的字符串
	signatureBase64 := base64.StdEncoding.EncodeToString(signature)

	// 打印签名
	//fmt.Println("Signature:", signatureBase64)
	return signatureBase64, nil
}

func (r ToolRsa) PrivateKeyToBytes() []byte {
	// 使用 x509 库将私钥编码为 DER 格式的字节切片
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(r.PrivateKey)

	// 使用 PEM 编码，方便存储和读取
	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	})

	//fmt.Println(privateKeyPEM)
	return privateKeyPEM

}

func (r *ToolRsa) PublicKeyToBytes() []byte {
	if r.PublicKey == nil || r.PublicKey.N == nil || r.PublicKey.E == 0 {
		return []byte("")
	}
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(r.PublicKey)
	if err != nil {
		return []byte("")
	}

	// 使用 PEM 编码，将公钥字节切片转换为字符串
	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY", // 公钥的类型
		Bytes: publicKeyBytes,
	})

	return publicKeyPEM

}

func (r *ToolRsa) RecoverRsaPublic(pemPublicKey []byte) *rsa.PublicKey {
	// 解码PEM格式的公钥
	block, _ := pem.Decode([]byte(pemPublicKey))
	if block == nil {
		log.Fatal("failed to parse PEM block containing the public key")
	}

	// 解析DER编码的公钥
	pubKeyInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		log.Fatal("failed to parse public key:", err)
	}

	// 将公钥转换为*rsa.PublicKey类型
	publicKey, ok := pubKeyInterface.(*rsa.PublicKey)
	if !ok {
		log.Fatal("failed to cast public key to *rsa.PublicKey")
	}
	r.PublicKey = publicKey

	return publicKey
}

func (r *ToolRsa) RecoverRsaPrivate(pemPrivateKey []byte) *rsa.PrivateKey {
	// 解码PEM格式的私钥
	block, _ := pem.Decode([]byte(pemPrivateKey))
	if block == nil {
		log.Fatal("failed to parse PEM block containing the private key")
	}

	// 解析DER编码的私钥
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		log.Fatal(err)
	}

	r.PrivateKey = privateKey
	return privateKey
}

func (r ToolRsa) FileCreateKey(key []byte, FileName string) error {
	// 保存到本地文件
	// 要检查的文件夹路径
	folderPath := "file/rsa/"

	// 检查文件夹是否存在
	_, err := os.Stat(folderPath)
	if os.IsNotExist(err) {
		// 文件夹不存在，创建它
		err := os.Mkdir(folderPath, 0755) // 0755 是文件夹权限
		if err != nil {
			fmt.Println("Failed to create folder:", err)
			return err
		}
		fmt.Println("Folder created successfully.")
	} else if err != nil {
		fmt.Println("Error:", err)
		return err
	} else {
		fmt.Println("Folder already exists.")
	}

	file, err := os.Create(FileName)
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer file.Close()

	_, err = file.WriteString(string(key))
	if err != nil {
		log.Fatal(err)
		return err
	}

	fmt.Println("RSA public key saved to " + FileName)
	return nil
}
