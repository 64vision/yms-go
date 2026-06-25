package main

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "wallet/proto"

	"google.golang.org/grpc"
)

const (
	wallerHost = "localhost:5051"
)

func main() {
	// Connect to gRPC server
	conn, err := grpc.Dial(wallerHost, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewWalletServiceClient(conn)

	// Call SayHello RPC
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	response, err := client.SendTransaction(ctx, &pb.TransactionRequest{
		AccountId:   801500000,
		Amount:      1000.501,
		Description: "Admin deposit",
		Type:        "deposit",
		SenderId:    516232,
	})
	if err != nil {
		log.Fatalf("Error calling OpenAccount: %v", err)
	}

	fmt.Println("Response from server:", response.Status)
	fmt.Println("Response from server:", response.Remarks)
	time.Sleep(10000)

	acct, err := client.CheckAccount(ctx, &pb.GetAccountRequest{
		AccountId: 801500000,
	})
	if err != nil {
		log.Fatalf("Error calling OpenAccount: %v", err)
	}

	fmt.Println("Response from server:", acct.Status)
	fmt.Println("Response from server:", acct.Remarks)
}
