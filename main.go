package main

import (
	"log"

	"github.com/micro/go-micro/v2"
	encrypt_proto "github.com/shunjiecloud-proto/encrypt/proto"
	"github.com/shunjiecloud/encrypt-srv/modules"
	"github.com/shunjiecloud/encrypt-srv/services"
)

func main() {
	//  Create srv
	service := micro.NewService(
		micro.Name("go.micro.srv.encrypt"),
	)

	//  init modules
	modules.Setup()

	//  init service
	service.Init()

	//  register Handlers
	encrypt_proto.RegisterEncryptHandler(service.Server(), new(services.EncryptService))

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
