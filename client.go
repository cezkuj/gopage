package main

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"

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
	go func() {
		for {
			_, err = client.Get1(context.Background(), &pb.GetParams{})
			if err != nil {
				log.Println(err)
				return
			}
		}
	}()

	for {
		_, err := client.Get2(context.Background(), &pb.GetParams{})
		if err != nil {
			log.Println(err)
			return
		}
	}
}
