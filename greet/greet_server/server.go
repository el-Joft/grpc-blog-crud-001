package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"time"

	"grpc-go-course-udemy/greet/greetpb"

	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
)

type server struct{}

func (*server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	fmt.Printf("Greet function was invoked %v\n", req)
	firstName := req.GetGreeting().GetFirstName()
	result := "Hello " + firstName
	res := &greetpb.GreetResponse{
		Result: result,
	}
	return res, nil
}

func (*server) GreetManyTimes(req *greetpb.GreetManyTimesRequest, stream greetpb.GreetService_GreetManyTimesServer) error {
	fmt.Printf("GreetManyTimes function invoked %v", req)
	firstName := req.GetGreet().GetFirstName()
	for i := 0; i < 10; i++ {
		result := "Hello " + firstName + " number " + strconv.Itoa(i)
		res := &greetpb.GreetManyTimeResponse{
			Result: result,
		}
		stream.Send(res)
		// time.Sleep(1000 * time.Millisecond)
	}
	return nil
}

func (*server) LongGreet(stream greetpb.GreetService_LongGreetServer) error {
	fmt.Printf("LongGreet function invoked with a streaming \n")

	result := "Hello "
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			// means we have reached the end of file
			return stream.SendAndClose(&greetpb.LongGreetResponse{
				Result: result,
			})
		}
		if err != nil {
			log.Fatalf("Error while reading client stream %v", err)

		}
		firstName := req.GetGreet().GetFirstName()
		result += "Hello " + firstName + "! "

	}
}

func (*server) GreetWithDealine(ctx context.Context, req *greetpb.GreetWithDeadlineRequest) (*greetpb.GreetWithDeadlineResponse, error) {
	fmt.Printf("GreetWithDeadline function was invoked %v\n", req)

	for i := 0; i < 3; i++ {
		if ctx.Err() == context.Canceled {
			// the client cancelled the request
			fmt.Println("The Client Cancelled the Request")
			return nil, status.Error(codes.Canceled, "The client cancelled the request")
		}
		time.Sleep(1 * time.Second)
	}
	firstName := req.GetGreet().GetFirstName()
	result := "Hello " + firstName
	res := &greetpb.GreetWithDeadlineResponse{
		Result: result,
	}

	return res, nil

}

func main() {
	fmt.Println("Server running")

	// pass in the tcp and the default address for rpc is 50051
	lis, err := net.Listen("tcp", "0.0.0.0:50051")

	if err != nil {
		log.Fatalf("Failed to Listen: %v", err)
	}

	keyFile := "ssl/server.pem"
	certFile := "ssl/server.crt"
	creds, sslErr := credentials.NewServerTLSFromFile(certFile, keyFile)
	if sslErr != nil {
		fmt.Println("Error Loading Certificate")
	}

	opts := []grpc.ServerOption{}
	opts = append(opts, grpc.Creds(creds))
	s := grpc.NewServer(opts...)
	greetpb.RegisterGreetServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to server %v", err)
	}

}
