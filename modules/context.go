package modules

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"

	"github.com/micro/go-micro/v2/config"
)

type moduleWrapper struct {
	DefaultKeyPair KeyPair
	RawPublicKey   string
}

//ModuleContext 模块上下文
var ModuleContext moduleWrapper

//Setup 初始化Modules
func Setup() {
	var defaultKeyPair KeyPairConfig
	if err := config.Get("config", "defaultKeyPair").Scan(&defaultKeyPair); err != nil {
		panic(err)
	}
	ModuleContext.RawPublicKey = defaultKeyPair.PublicKey
	//  publicKey
	block, _ := pem.Decode([]byte(defaultKeyPair.PublicKey)) //将密钥解析成公钥实例
	if block == nil {
		panic("public key decode failed")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes) //解析pem.Decode（）返回的Block指针实例
	if err != nil {
		panic(err)
	}
	isOk := false
	ModuleContext.DefaultKeyPair.PublicKey, isOk = pubInterface.(*rsa.PublicKey)
	if isOk == false {
		panic("public interface convert to publicKey failed")
	}

	//  privateKey
	block, _ = pem.Decode([]byte(defaultKeyPair.PrivateKey)) //将密钥解析成私钥实例
	if block == nil {
		panic("private key decode failed")
	}
	ModuleContext.DefaultKeyPair.PrivateKey, err = x509.ParsePKCS1PrivateKey(block.Bytes) //解析pem.Decode（）返回的Block指针实例
	if err != nil {
		panic("private key parse failed")
	}
}
