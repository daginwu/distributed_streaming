package api

import (
	"context"
	pb "distributed_streaming/cmd/persist/app/datatype/pb"
	"log"
	"strconv"

	badger "github.com/dgraph-io/badger/v3"
	"go.uber.org/fx"
	"google.golang.org/grpc"
)

var Modual = fx.Options(

	fx.Provide(),
	fx.Invoke(
		InitAPI,
	),
)

type server struct {
	pb.UnimplementedPersistServiceServer
}

var db *badger.DB

func (s *server) CreateUser(ctx context.Context, in *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	log.Println("[CreateUser]")
	log.Println("id: " + in.Id)
	log.Println("name: " + in.Name)

	// Persist to badger

	err := db.Update(func(txn *badger.Txn) error {
		txn.Set(
			[]byte("users:"+in.Id+":name"),
			[]byte(in.Name),
		)
		txn.Set(
			[]byte("users:"+in.Id+":balance"),
			[]byte(strconv.Itoa(int(in.Balance))),
		)
		return nil
	})

	if err != nil {
		return &pb.CreateUserResponse{
			Reply: "Create user fail",
		}, nil
	}
	return &pb.CreateUserResponse{
		Reply: "Create user with id: " + in.Id,
	}, nil
}

func InitAPI(s *grpc.Server, database *badger.DB) {
	pb.RegisterPersistServiceServer(s, &server{})
	db = database
}
