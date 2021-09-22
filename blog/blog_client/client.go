package main

import (
	"context"
	"fmt"
	"grpc-go-course-udemy/blog/blogpb"
	"log"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Blog Client Running")
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}

	defer cc.Close()

	// create the client
	c := blogpb.NewBlogServiceClient(cc)

	fmt.Println("Creating Blog")
	blog := &blogpb.Blog{
		AuthorId: "Stephane",
		Title:    "Title of my First gRPC Blog",
		Content:  "Content of my first gRPC blog",
	}
	res, err := c.CreateBlog(context.Background(), &blogpb.CreateBlogRequest{
		Blog: blog,
	})
	if err != nil {
		log.Fatalf("Unexpected Error %v", err)
	}
	fmt.Printf("Blog has been created: %v", res)
}
