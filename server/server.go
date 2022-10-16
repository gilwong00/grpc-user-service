package main

import (
	"context"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"os"

	pb "github.com/gilwong00/grpc-user-service/user"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
)

const (
	port = ":8080"
)

type UserService struct {
	pb.UnimplementedUserServiceServer
	// user_list *pb.UserList
}

func writeUserToFile(users *pb.UserList) error {
	jsonBytes, err := protojson.Marshal(users)

	if err != nil {
		log.Fatalf("JSON Marshalling failed %v", err)
		return err
	}

	if err := ioutil.WriteFile("users.json", jsonBytes, 0664); err != nil {
		return err
	}

	return nil
}

func (s *UserService) CreateNewUser(ctx context.Context, in *pb.NewUser) (*pb.User, error) {
	log.Printf("Received: %v", in.GetName())
	readBytes, err := ioutil.ReadFile("users.json")
	var users_list *pb.UserList = &pb.UserList{}
	var user_id int32 = int32(rand.Intn(1000))
	created_user := &pb.User{
		Name: in.GetName(), Age: in.GetAge(), Id: user_id,
	}

	// create users.json for the first time if file doesnt exist
	if err != nil {
		if os.IsNotExist(err) {
			log.Print("File not found. Creating a new file")
			users_list.Users = append(users_list.Users, created_user)

			if err := writeUserToFile(users_list); err != nil {
				log.Fatalf("Failed writing to file %v", err)
			}

			return created_user, nil
		} else {
			log.Fatalln("Error reading file: ", err)
		}
	}

	if err := protojson.Unmarshal(readBytes, users_list); err != nil {
		log.Fatalf("Failed to parse user list: %v", err)
	}

	users_list.Users = append(users_list.Users, created_user)

	if err := writeUserToFile(users_list); err != nil {
		log.Fatalf("Failed writing to file %v", err)
	}

	// s.user_list.Users = append(s.user_list.Users, created_user)
	return created_user, nil
}

func (s *UserService) GetUsers(ctx context.Context, in *pb.GetUsersParams) (*pb.UserList, error) {
	// return s.user_list, nil
	jsonBytes, err := ioutil.ReadFile("users.json")
	if err != nil {
		log.Fatalf("Failed to read from file: %v", err)
	}

	var users_list *pb.UserList = &pb.UserList{}
	if err := protojson.Unmarshal(jsonBytes, users_list); err != nil {
		log.Fatalf("Unmarshalling failed: %v", err)
	}
	return users_list, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterUserServiceServer(s, &UserService{
		// user_list: &pb.UserList{},
	})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
