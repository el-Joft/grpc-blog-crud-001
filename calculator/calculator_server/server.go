package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"math"
	"net"

	"grpc-go-course-udemy/calculator/calculatorpb"

	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/reflection"
	
)

type server struct{}

func (*server) Sum(ctx context.Context, req *calculatorpb.SumRequest) (*calculatorpb.SumResponse, error) {
	fmt.Printf("Received Sum RPC: %v\n ", req)
	firstNumber := req.FirstNumber
	secondNumber := req.SecondNumber
	result := firstNumber + secondNumber
	res := &calculatorpb.SumResponse{
		SumResult: result,
	}
	return res, nil
}

func (*server) ComputeAverage(stream calculatorpb.CalculatorService_ComputeAverageServer) error {
	fmt.Printf("CalculateAverage function invoked with a streaming \n")

	var result int32 = 0
	var counter int32 = 0
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			average := float64(result) / float64(counter)
			return stream.SendAndClose(&calculatorpb.ComputeAverageResponse{
				Result: int32(average),
			})
		}
		if err != nil {
			log.Fatalf("Error while reading client stream %v", err)
		}
		parsedNumber := req.GetNumber()
		result += parsedNumber
		counter++
	}
}

func (*server) SquareRoot(ctx context.Context, req *calculatorpb.SquareRootRequest) (*calculatorpb.SquareRootResponse, error) {
	fmt.Printf("Recieved Square Root RPC: %v\n", req)
	number := req.GetNumber()
	if number < 0 {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Recieved a negative number: %v\n", number),
		)
	}
	return &calculatorpb.SquareRootResponse{
		NumberRoot: math.Sqrt(float64(number)),
	}, nil
}

func main() {
	fmt.Println("Calculator Service is running")

	// pass in the tcp and the default address for rpc is 50051
	lis, err := net.Listen("tcp", "0.0.0.0:50051")

	if err != nil {
		log.Fatalf("Failed to Listen: %v", err)
	}

	s := grpc.NewServer()
	calculatorpb.RegisterCalculatorServiceServer(s, &server{})

	// Register reflection service on gRPC server.
	reflection.Register(s)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to server %v", err)
	}

}
