package transport

import (
	"context"
	"fmt"

	//"context"
	//"fmt"
	"google.golang.org/grpc"
	"image-gallery/internal/gallery/config"
	pb "image-gallery/pkg/protobuf/authorizationservice/gw"
)

type AuthGrpcTransport struct {
	config config.AuthGrpcTransport
	client pb.AuthorizationServiceClient
}

func NewAuthGrpcTransport(config config.AuthGrpcTransport) *AuthGrpcTransport {
	opts := []grpc.DialOption{grpc.WithInsecure()}

	conn, _ := grpc.Dial(config.Host, opts...)

	client := pb.NewAuthorizationServiceClient(conn)

	return &AuthGrpcTransport{
		client: client,
		config: config,
	}
}

func (t *AuthGrpcTransport) IsUserAuthorized(ctx context.Context, tokenString string) (*pb.UserAuthorizationResponse, error) {
	resp, err := t.client.IsUserAuthorized(ctx, &pb.UserAuthorizationRequest{
		TokenString: tokenString,
	})

	if err != nil {
		return nil, fmt.Errorf("Can not check Authorizaion: %s", err)
	}
	if resp == nil {
		return nil, fmt.Errorf("not found")
	}

	return resp, nil
}
