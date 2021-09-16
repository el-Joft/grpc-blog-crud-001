package main

import (
	"context"
	"fmt"
	"grpc-go-course-udemy/calculator/calculatorpb"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func main() {
	fmt.Println("Calculator Client Running")

	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}

	defer cc.Close()

	// create the client
	c := calculatorpb.NewCalculatorServiceClient(cc)

	// doUnary(c)
	// doClientStreaming(c)
	doErrorUnary(c)

}

func doUnary(c calculatorpb.CalculatorServiceClient) {
	fmt.Printf("Starting to do a Unary RPC")
	req := &calculatorpb.SumRequest{
		FirstNumber:  12,
		SecondNumber: 40,
	}

	res, err := c.Sum(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling Sum RPC %v ", err)
	}
	log.Printf("Response from Sum: %v", res)

}

func doClientStreaming(c calculatorpb.CalculatorServiceClient) {
	fmt.Printf("Starting to do a server streaming RPC")

	// requests := []*calculatorpb.ComputeAverageRequest{
	// 	&calculatorpb.ComputeAverageRequest{
	// 		Number: 1,
	// 	},
	// 	&calculatorpb.ComputeAverageRequest{
	// 		Number: 2,
	// 	},
	// 	&calculatorpb.ComputeAverageRequest{
	// 		Number: 3,
	// 	},
	// 	&calculatorpb.ComputeAverageRequest{
	// 		Number: 4,
	// 	},
	// }
	stream, err := c.ComputeAverage(context.Background())
	if err != nil {
		log.Fatalf("Error while calling Compute Calculate Client: %v", err)
	}

	numbers := []int32{5, 3, 24, 65}

	for _, number := range numbers {
		stream.Send(
			&calculatorpb.ComputeAverageRequest{
				Number: number,
			},
		)

		time.Sleep(1000 * time.Millisecond)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error while recieving response from Compute Server %v", err)

	}
	fmt.Printf("Calculate Response:  %v ", res)
}

func doErrorUnary(c calculatorpb.CalculatorServiceClient) {
	fmt.Printf("Starting to to a Unary Error Streaming RPC")

	// correct call
	var number int32 = 20
	doErrorCall(c, number)

	// invalid call
	var invalidNumber int32 = -2
	doErrorCall(c, invalidNumber)

}

func doErrorCall(c calculatorpb.CalculatorServiceClient, number int32) {
	res, err := c.SquareRoot(context.Background(), &calculatorpb.SquareRootRequest{
		Number: number,
	})
	if err != nil {
		respErr, ok := status.FromError(err)
		if ok {
			// Meaning it is actuall error from gRPC
			fmt.Printf("Error Message from Server %v \n", respErr.Message())
			fmt.Println(respErr.Code())
			if respErr.Code() == codes.InvalidArgument {
				fmt.Println("We probaly sent a negative number")
				return
			}
		} else {
			log.Fatalf("Big Error calling SquareRoot %v \n", err)
			return 
		}
	}
	fmt.Printf("Result of square root of %v = %v\n", number, res.GetNumberRoot())
}
