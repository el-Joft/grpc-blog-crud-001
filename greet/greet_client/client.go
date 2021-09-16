package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"grpc-go-course-udemy/greet/greetpb"

	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
)

func main() {
	fmt.Println("Hello I am a client")

	// create a connection to the server
	// we are using WithInsecure because gRPC comes with SSL by default
	// cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	// using ssl
	certFile := "ssl/ca.crt" // Certificate Authority Trust Certificate
	creds, certErr := credentials.NewClientTLSFromFile(certFile, "")
	if certErr != nil {
		fmt.Printf("Error while loading certfile %v \n", certErr)
	}

	opts := grpc.WithTransportCredentials(creds)
	conn, err := grpc.Dial("localhost:50051", opts)
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}

	// defer is usually called at the end of the code
	defer conn.Close()

	// create a client
	c := greetpb.NewGreetServiceClient(conn)
	doUnary(c)
	// doServerStreaming(c)
	// doClientStreaming(c)
	// doUnaryWithDeadline(c, 5*time.Second) // Should complete
	// doUnaryWithDeadline(c, 1*time.Second) // Should timeout
}

func doUnary(c greetpb.GreetServiceClient) {
	fmt.Printf("Starting to do a Unary RPC")
	// fmt.Printf("Created Client %v", c)
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Stephane",
			LastName:  "Maleek",
		},
	}
	res, err := c.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling Greet RPC: %v", err)
	}
	log.Printf("Response from Greet: %v", res)
}

func doServerStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a server streaming RPC")
	req := &greetpb.GreetManyTimesRequest{
		Greet: &greetpb.Greeting{
			FirstName: "Stephane",
			LastName:  "Maleek",
		},
	}
	resStream, err := c.GreetManyTimes(context.Background(), req)

	if err != nil {
		log.Fatalf("error while calling GreetManyTimes RPC: %v", err)
	}

	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			// we have reached the end of the stream
			break
		}
		if err != nil {
			log.Fatalf("Error while reading stream: %v", err)
		}
		log.Printf("Response from Greet: %v", msg.GetResult())
	}

}

func doClientStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a Client streaming RPC")

	requests := []*greetpb.LongGreetRequest{
		&greetpb.LongGreetRequest{
			Greet: &greetpb.Greeting{
				FirstName: "Stephane",
			},
		},
		&greetpb.LongGreetRequest{
			Greet: &greetpb.Greeting{
				FirstName: "John",
			},
		},
		&greetpb.LongGreetRequest{
			Greet: &greetpb.Greeting{
				FirstName: "Lucy",
			},
		},
		&greetpb.LongGreetRequest{
			Greet: &greetpb.Greeting{
				FirstName: "Mark",
			},
		},
		&greetpb.LongGreetRequest{
			Greet: &greetpb.Greeting{
				FirstName: "Piper",
			},
		},
	}

	stream, err := c.LongGreet(context.Background())
	if err != nil {
		log.Fatalf("Error while calling Long greet %v", err)
	}

	// we iterate over our slice and send the message individually
	for _, req := range requests {
		fmt.Printf("Sending request: %v \n", req)
		stream.Send(req)
		time.Sleep(1000 * time.Millisecond)
	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error while recieving response from LongGreet %v", err)
	}
	fmt.Printf("LongGreet Response: %v \n ", res)
}

func doUnaryWithDeadline(c greetpb.GreetServiceClient, timeout time.Duration) {
	fmt.Printf("Starting to do a Unary With Deadline RPC")
	// fmt.Printf("Created Client %v", c)
	req := &greetpb.GreetWithDeadlineRequest{
		Greet: &greetpb.Greeting{
			FirstName: "Stephane",
			LastName:  "Maleek",
		},
	}
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	// when the funtion completes we call defer
	defer cancel()
	res, err := c.GreetWithDealine(ctx, req)
	if err != nil {
		statusErr, ok := status.FromError(err)
		if ok {
			// Meaning it is actuall error from gRPC
			fmt.Printf("Error Message from Server %v \n", statusErr.Message())
			fmt.Println(statusErr.Code())
			if statusErr.Code() == codes.DeadlineExceeded {
				fmt.Println("Timeout was hit! Deadline was exceeded")

			} else {
				fmt.Printf("Unexpected Error: %v", err)
			}
		} else {
			log.Fatalf("Big Error calling GreetWithDealine %v \n", err)
			return
		}
	}
	log.Printf("Response from Greet: %v", res)
}
