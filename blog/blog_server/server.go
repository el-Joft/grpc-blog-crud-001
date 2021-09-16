package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"grpc-go-course-udemy/blog/blogpb"

	"google.golang.org/grpc"
)

type server struct {
}

func main() {
	// if we crash the go code we get the file name and line number
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	fmt.Println("Blog Server running")

	// pass in the tcp and the default address for rpc is 50051
	lis, err := net.Listen("tcp", "0.0.0.0:50051")

	if err != nil {
		log.Fatalf("Failed to Listen: %v", err)
	}

	opts := []grpc.ServerOption{}
	s := grpc.NewServer(opts...)
	blogpb.RegisterBlogServiceServer(s, &server{})

	// Gracefully shutdown the server
	go func() {
		fmt.Println("Starting Server....")
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Failed to server %v", err)
		}
	}()

	// wait for control-c to exit
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	// block until signal is received
	<-ch
	fmt.Println("Stopping the Server")
	s.Stop()
	fmt.Println("Closing the Listener")
	lis.Close()
	fmt.Println("End of Program")
}
