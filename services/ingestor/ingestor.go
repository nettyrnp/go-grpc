package main

import (
	"log"
	"time"

	"bufio"
	"encoding/csv"
	pb "github.com/nettyrnp/go-grpc/proto"
	"github.com/nettyrnp/go-grpc/util"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"os"
	"os/signal"
	"syscall"
)

const (
	address = "localhost:50502" // TODO: Move to config file
)

func main() {
	// Handle termination signals
	errors := make(chan error, 1)
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	signal.Notify(signals, syscall.SIGTERM)

	go func() {
		fname := "db/data.csv" // TODO: Move to config file
		var offset = 0
		limit := 40 // max lines to load

		csvFile, err := os.Open(fname)
		if err != nil {
			log.Fatalf("Error: %v", err)
		}
		defer csvFile.Close()
		bufreader := bufio.NewReader(csvFile)
		reader := csv.NewReader(bufreader)

		// Set up a connection to the server
		conn, err := grpc.Dial(address, grpc.WithInsecure())
		if err != nil {
			log.Fatalf("Did not connect: %v", err)
		}
		defer conn.Close()
		c := pb.NewPersistorClient(conn)

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		var err2 error
		for {
			// Read a portion of file
			people := util.Ingest(reader, offset, limit)
			if len(people) == 0 {
				break
			}
			log.Printf("Reading: loaded %d non-duplicate lines", len(people))

			// Save a portion of file
			r2, err2 := c.SavePersons(ctx, &pb.PersonsReq{
				Persons: fromModel(people),
			})
			if err2 != nil {
				log.Fatalf("Failed to save records: %v", err2)
				break
			}
			log.Printf("Saving: created %d, updated %d records in DB", r2.CreatedCount, r2.UpdatedCount)
			offset = +limit

			// Slow down to better see the effect of <Ctrl+C>
			time.Sleep(500 * time.Millisecond)
		}
		if err != nil {
			errors <- err
		}
		if err2 != nil {
			//log.Fatalf("client failed: %v", err2)
			errors <- err2
		}
	}()

	// Wait until the service fails or it is terminated.
	select {
	case err := <-errors:
		// Handle the error from DigestAndPersist
		log.Printf("Error from DigestAndPersist: %v\n", err)
		break
	case sig := <-signals:
		// Handle shutdown signals
		log.Printf("Signal: %v\n", sig)
		break
	}

	// Gracefully terminate
	i := 3
	for i > 0 {
		log.Printf("Terminating ingestor service in %d s\n", i)
		time.Sleep(1 * time.Second)
		i = i - 1
	}
}

func fromModel(people []util.Person) []*pb.Person {
	var people2 []*pb.Person
	for _, p := range people {
		people2 = append(people2, &pb.Person{
			Id:           p.Id,
			Name:         p.Name,
			Email:        p.Email,
			MobileNumber: p.MobileNumber,
		})
	}
	return people2
}
