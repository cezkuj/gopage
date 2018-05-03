package main

import (
"log"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	pb "github.com/cezkuj/gopage-feed/grpc"
)

func startClient() {
	conn, err := grpc.Dial("localhost:8002", grpc.WithInsecure())
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()
	client := pb.NewGopageClient(conn)
	number, err := client.Get1(context.Background(), &pb.GetParams{})
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(number)
}
