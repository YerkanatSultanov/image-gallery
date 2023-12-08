package transport

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"image-gallery/internal/auth/config"
	pb "image-gallery/pkg/protobuf/userservice/gw"
)

type UserGrpcTransport struct {
	config config.UserGrpcTransport
	client pb.UserServiceClient
}

func NewUserGrpcTransport(config config.UserGrpcTransport) *UserGrpcTransport {
	//nolint:all
	opts := []grpc.DialOption{grpc.WithInsecure()}

	conn, _ := grpc.Dial(config.Host, opts...)

	client := pb.NewUserServiceClient(conn)

	return &UserGrpcTransport{
		client: client,
		config: config,
	}
}

func (t *UserGrpcTransport) GetUserByEmail(ctx context.Context, email string) (*pb.User, error) {

	resp, err := t.client.GetUserByEmail(ctx, &pb.GetUserByEmailRequest{
		Email: email,
	})

	if err != nil {
		return nil, fmt.Errorf("cannot GetUserByEmail: %w", err)
	}

	if resp == nil {
		return nil, fmt.Errorf("not found")
	}

	return resp.Result, nil
}
