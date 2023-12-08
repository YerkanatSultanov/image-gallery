package transport

import (
	"context"
	"fmt"

	//"context"
	//"fmt"
	"google.golang.org/grpc"
	"image-gallery/internal/gallery/config"
	pb "image-gallery/pkg/protobuf/userservice/gw"
)

type UserGrpc struct {
	config config.UserGrpcTransport
	client pb.UserServiceClient
}

func NewUserGrpcTransport(config config.UserGrpcTransport) *UserGrpc {
	//nolint:all
	opts := []grpc.DialOption{grpc.WithInsecure()}

	conn, _ := grpc.Dial(config.Host, opts...)

	client := pb.NewUserServiceClient(conn)

	return &UserGrpc{
		client: client,
		config: config,
	}
}

func (t *UserGrpc) GetUserById(ctx context.Context, id int) (*pb.User, error) {
	resp, err := t.client.GetUserById(ctx, &pb.GetUserByIdRequest{
		Id: int32(id),
	})

	if err != nil {
		return nil, fmt.Errorf("Can not check Authorizaion: %s", err)
	}
	if resp == nil {
		return nil, fmt.Errorf("not found")
	}

	return resp.Result, nil
}

func (t *UserGrpc) GetUserByUsername(ctx context.Context, username string) (*pb.User, error) {
	resp, err := t.client.GetUserByUsername(ctx, &pb.GetUserByUsernameRequest{
		Username: username,
	})

	if err != nil {
		return nil, fmt.Errorf("Can not check Authorizaion: %s", err)
	}
	if resp == nil {
		return nil, fmt.Errorf("not found")
	}

	return resp.Result, nil
}
