package services

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"

	proto "github.com/shunjiecloud-proto/encrypt/proto"
	"github.com/shunjiecloud/encrypt-srv/modules"
)

type EncryptService struct{}

func (h *EncryptService) GetPublicKey(ctx context.Context, in *proto.GetPublicKeyRequest, out *proto.GetPublicKeyResponse) error {
	out.PublicKey = modules.ModuleContext.RawPublicKey
	return nil
}

func (h *EncryptService) Encrypt(ctx context.Context, in *proto.EncryptRequest, out *proto.EncryptResponse) error {
	ciphertext, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, modules.ModuleContext.DefaultKeyPair.PublicKey, []byte(in.Original), nil)
	if err != nil {
		return err
	}
	//  base64
	base64CipherText := base64.StdEncoding.EncodeToString(ciphertext)
	out.CipherText = base64CipherText
	return nil
}

func (h *EncryptService) Decrypt(ctx context.Context, in *proto.DecryptRequest, out *proto.DecryptResponse) error {
	//  base64 decode
	cipherText, err := base64.StdEncoding.DecodeString(in.CipherText)
	if err != nil {
		return err
	}
	original, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, modules.ModuleContext.DefaultKeyPair.PrivateKey, cipherText, nil)
	if err != nil {
		return err
	}
	out.Original = string(original)
	return nil
}
