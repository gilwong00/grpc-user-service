package main

import (
	"context"
	"log"
	"math/rand"
	"net"

	pb "github.com/gilwong00/grpc-user-service/user"

	"google.golang.org/grpc"
)

const (
	port = ":8080"
)

type UserService struct {
	pb.UnimplementedUserServiceServer
}

func (s *UserService) CreateNewUser(ctx context.Context, in *pb.NewUser) (*pb.User, error) {
	log.Printf("Received: %v", in.GetName())
	var user_id int32 = int32(rand.Intn(1000))
	return &pb.User{
		Name: in.GetName(), Age: in.GetAge(), Id: user_id,
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterUserServiceServer(s, &UserService{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
