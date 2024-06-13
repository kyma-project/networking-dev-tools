package hello_server

import (
	"context"
	"github.com/kyma-project/networking-dev-tools/grpcbin/pkg/hello"
)

type Server struct {
	hello.UnimplementedHelloServiceServer
}

func (h Server) SayHello(_ context.Context, request *hello.HelloRequest) (*hello.HelloResponse, error) {
	return &hello.HelloResponse{Message: "Hello " + request.Name}, nil
}

func (h Server) StreamGoats(_ *hello.Pen, stream hello.HelloService_StreamGoatsServer) error {
	goats := []*hello.Goat{
		{Name: "A"},
		{Name: "B"},
		{Name: "C"},
	}

	for _, goat := range goats {
		if err := stream.Send(goat); err != nil {
			return err
		}
	}

	return nil
}

func (h Server) ListGoats(_ context.Context, _ *hello.Pen) (*hello.GoatList, error) {
	return &hello.GoatList{
		Goats: []*hello.Goat{
			{Name: "A"},
			{Name: "B"},
			{Name: "C"},
		},
	}, nil
}
