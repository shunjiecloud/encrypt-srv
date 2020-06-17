module github.com/shunjiecloud/encrypt-srv

go 1.14

require (
	github.com/go-redis/redis/v7 v7.4.0
	github.com/micro/go-micro/v2 v2.8.0
	github.com/shunjiecloud-proto/encrypt v0.0.0-20200617090148-a99de7331466
	github.com/shunjiecloud/errors v1.0.3-0.20200427091440-d2c8251bbc81
	github.com/shunjiecloud/pkg v0.0.0-20200608213205-7936a725a0c8
	golang.org/x/net v0.0.0-20200602114024-627f9648deb9 // indirect
	golang.org/x/sys v0.0.0-20200602225109-6fdc65e7d980 // indirect
	google.golang.org/protobuf v1.24.0 // indirect
)

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0
