package main

import (
	"log"
	"net"

	"fmt"
	"github.com/nettyrnp/go-grpc/db"
	pb "github.com/nettyrnp/go-grpc/proto"
	"github.com/nettyrnp/go-grpc/util"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	port = ":50502"
)

type server struct{}

func (s *server) SavePersons(ctx context.Context, in *pb.PersonsReq) (*pb.PersonsReply, error) {
	repl := &pb.PersonsReply{
		CreatedCount: 0,
		UpdatedCount: 0,
	}
	db.ResetDB()
	people := fromProto(in.Persons)
	for _, p := range people {
		//fmt.Println("received p:", p)
		c1, c2, err := db.SaveRecord(p)
		if err != nil {
			panic(fmt.Errorf("Error while saving record to database: %s", err))
		}
		repl.CreatedCount += c1
		repl.UpdatedCount += c2
	}
	return repl, nil
}

func fromProto(people []*pb.Person) []util.Person {
	var people2 []util.Person
	for _, p := range people {
		//fmt.Println("p:", p.ToString())
		people2 = append(people2, util.Person{
			Id:           p.Id,
			Name:         p.Name,
			Email:        p.Email,
			MobileNumber: p.MobileNumber,
		})
	}
	return people2
}

func main() {
	// Handle termination signals
	errors := make(chan error, 1)
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	signal.Notify(signals, syscall.SIGTERM)

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterPersistorServer(s, &server{})
	// Register reflection service on gRPC server.
	reflection.Register(s)

	go func() {
		err := s.Serve(lis)
		if err != nil {
			log.Fatalf("Failed to serve: %v", err)
			errors <- err
		}
	}()

	// Wait until the service fails or it is terminated.
	select {
	case err := <-errors:
		// Handle the error from s.Serve
		log.Printf("Error from s.Serve: %v\n", err)
		break
	case sig := <-signals:
		// Handle shutdown signals
		log.Printf("Signal: %v\n", sig)
		break
	}

	// Gracefully terminate
	i := 3
	for i > 0 {
		log.Printf("Terminating persistence service in %d s\n", i)
		time.Sleep(1 * time.Second)
		i = i - 1
	}
}
