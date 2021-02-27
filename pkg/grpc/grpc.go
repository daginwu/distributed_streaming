package grpc

import (
	"context"
	"net"

	"go.uber.org/fx"
	"google.golang.org/grpc"

	"log"

	"github.com/spf13/viper"
)

var Modual = fx.Options(

	fx.Provide(
		NewgGrpcServer,
	),
	fx.Invoke(
		InitGrpcServer,
	),
)

func NewgGrpcServer(lc fx.Lifecycle) *grpc.Server {

	log.Println("[Distributed_Streaming] gRPC server start")

	port := viper.GetString("grpc.port")

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()

	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			go func() error {
				if err := s.Serve(lis); err != nil {
					return err
				}
				return nil
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			s.GracefulStop()
			return nil
		},
	})

	return s
}

func InitGrpcServer(s *grpc.Server) error {
	return nil
}
