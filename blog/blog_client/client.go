package main

import (
	"context"
	"fmt"
	"grpc-go-course-udemy/blog/blogpb"
	"io"
	"log"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Blog Client Running")
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Printf("Could not connect: %v", err)
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
		log.Printf("Unexpected Error %v \n", err)
	}
	fmt.Printf("Blog has been created: %v \n", res)

	// Read a Blog
	fmt.Println("Reading a Blog")
	_, readBlogErr := c.ReadBlog(context.Background(), &blogpb.ReadBlogRequest{
		BlogId: "2382392832932",
	})
	if readBlogErr != nil {
		log.Printf("Unexpected Error while Reading Blog %v", readBlogErr)
	}

	readBlogWithoutError, readBlogErrWithoutError := c.ReadBlog(context.Background(), &blogpb.ReadBlogRequest{
		BlogId: res.GetBlog().Id,
	})
	if readBlogErr != nil {
		log.Printf("Unexpected Error while Reading Blog %v \n", readBlogErrWithoutError)
	}
	fmt.Printf("Blog Read Successfully %v \n", readBlogWithoutError)

	// Create another blog
	newBlog := &blogpb.Blog{
		Id:       res.GetBlog().Id,
		AuthorId: "Changed Author",
		Title:    "Title of the Changed gRPC Blog",
		Content:  "Content of the Changed gRPC blog",
	}
	updateBlogRes, updateBlogErr := c.UpdateBlog(context.Background(), &blogpb.UpdateBlogRequest{
		Blog: newBlog,
	})
	if err != nil {
		log.Printf("Unexpected Error %v \n", updateBlogErr)
	}
	fmt.Printf("Blog updated Successfully %v \n", updateBlogRes)

	// Delete Blog
	deleteBlogRes, deleteBlogError := c.DeleteBlog(context.Background(), &blogpb.DeleteBlogRequest{
		BlogId: res.GetBlog().Id,
	})
	if err != nil {
		log.Printf("Unexpected Error %v \n", deleteBlogError)
	}
	fmt.Printf("Blog deleted Successfully %v \n", deleteBlogRes)

	// List Blog
	stream, listBlogErr := c.ListBlog(context.Background(), &blogpb.ListBlogRequest{})
	if listBlogErr != nil {
		log.Fatalf("Error while calling ListBlog RPC %v", listBlogErr)
	}

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Something Happened: %v", err)
		}
		fmt.Println(res.GetBlog())
	}

}
