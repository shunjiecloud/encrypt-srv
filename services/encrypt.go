package services

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"time"

	"github.com/go-redis/redis/v7"
	merr "github.com/micro/go-micro/v2/errors"
	proto "github.com/shunjiecloud-proto/encrypt/proto"
	"github.com/shunjiecloud/encrypt-srv/modules"
	"github.com/shunjiecloud/errors"
)

type EncryptService struct{}

const KeyPairIdSequence = "encrypt:keypair:sequence"
const KeyPairRedisKey = "encrypt:keypair:%v"

type keyPair struct {
	PublicKey  []byte `json:"public_key"`
	PrivateKey []byte `json:"private_key"`
}

func setKeyPair(redis *redis.Client, publicKey, privateKey []byte) (id int64, err error) {
	id, err = redis.Incr(KeyPairIdSequence).Result()
	if err != nil {
		return id, err
	}
	key := fmt.Sprintf(KeyPairRedisKey, id)
	keyPair := keyPair{
		PublicKey:  publicKey,
		PrivateKey: privateKey,
	}
	j, err := json.Marshal(&keyPair)
	if err != nil {
		return id, err
	}
	_, err = redis.Set(key, j, time.Duration(5)*time.Minute).Result()
	if err != nil {
		return id, err
	}
	return id, nil
}

func getKeyPair(redis *redis.Client, id int64) (*keyPair, error) {
	key := fmt.Sprintf(KeyPairRedisKey, id)
	j, err := redis.Get(key).Result()
	if err != nil {
		return nil, err
	}
	var keyPair keyPair
	err = json.Unmarshal([]byte(j), &keyPair)
	if err != nil {
		return nil, err
	}
	return &keyPair, nil
}

func (h *EncryptService) GetPublicKey(ctx context.Context, in *proto.GetPublicKeyRequest, out *proto.GetPublicKeyResponse) error {
	//  生成密钥对
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}
	publicKey := &privateKey.PublicKey

	//  生成私钥
	PKCS1PrivateKey := x509.MarshalPKCS1PrivateKey(privateKey)
	if err != nil {
		return err
	}
	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: PKCS1PrivateKey,
	}
	buffPrivateKey := bytes.NewBuffer(nil)
	err = pem.Encode(buffPrivateKey, block)
	if err != nil {
		return err
	}

	//  生成公钥
	PKIXPublicKey, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return err
	}
	block = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: PKIXPublicKey,
	}
	buffPublicKey := bytes.NewBuffer(nil)
	err = pem.Encode(buffPublicKey, block)
	if err != nil {
		return err
	}

	//  保存到redis
	id, err := setKeyPair(modules.ModuleContext.Redis, buffPublicKey.Bytes(), buffPrivateKey.Bytes())
	if err != nil {
		return err
	}
	//  返回
	out.PublicKey = base64.StdEncoding.EncodeToString(buffPublicKey.Bytes())
	out.PublicKeyId = id
	return nil
}

func (h *EncryptService) Encrypt(ctx context.Context, in *proto.EncryptRequest, out *proto.EncryptResponse) error {
	//  取得密钥对
	keyPair, err := getKeyPair(modules.ModuleContext.Redis, in.PublicKeyId)
	if err != nil {
		return err
	}
	//  publicKey
	block, _ := pem.Decode(keyPair.PublicKey) //将密钥解析成公钥实例
	if block == nil {
		return err
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes) //解析pem.Decode（）返回的Block指针实例
	if err != nil {
		panic(err)
	}
	publicKey, isOk := pubInterface.(*rsa.PublicKey)
	if isOk == false {
		return errors.New(merr.BadRequest("get publicKey failed", ""))
	}
	ciphertext, err := rsa.EncryptOAEP(sha1.New(), rand.Reader, publicKey, []byte(in.Original), nil)
	if err != nil {
		return err
	}
	//  base64
	base64CipherText := base64.StdEncoding.EncodeToString(ciphertext)
	out.CipherText = base64CipherText
	return nil
}

func (h *EncryptService) Decrypt(ctx context.Context, in *proto.DecryptRequest, out *proto.DecryptResponse) error {
	//  取得密钥对
	keyPair, err := getKeyPair(modules.ModuleContext.Redis, in.PublicKeyId)
	if err != nil {
		return err
	}
	block, _ := pem.Decode(keyPair.PrivateKey) //将密钥解析成私钥实例
	if block == nil {
		return errors.New(merr.BadRequest("private key decode failed", ""))
	}
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes) //解析pem.Decode（）返回的Block指针实例
	//  base64 decode
	cipherText, err := base64.StdEncoding.DecodeString(in.CipherText)
	if err != nil {
		return err
	}
	original, err := rsa.DecryptOAEP(sha1.New(), rand.Reader, privateKey, cipherText, nil)
	if err != nil {
		return err
	}
	out.Original = string(original)
	return nil
}
