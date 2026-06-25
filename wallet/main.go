package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/go-pg/pg"
	"google.golang.org/grpc"
	pb "hyperball.com/wallet/proto"
)

var DBM *pg.DB

const (
	PORT = ":5051"
)

// HelloServiceServer implements the gRPC service
type WalletServiceServer struct {
	pb.UnimplementedWalletServiceServer
}

func (s *WalletServiceServer) CheckAccount(ctx context.Context, req *pb.GetAccountRequest) (*pb.BalanceResponse, error) {
	fmt.Println("Check wallet account...")

	acct := InquireAccount(int(req.AccountId))
	if !acct.Status {
		return &pb.BalanceResponse{Remarks: acct.Remarks, AccountId: int32(acct.AccountID), Status: false}, nil
	}
	return &pb.BalanceResponse{Remarks: acct.Remarks, AccountId: int32(acct.AccountID), Status: true}, nil

}

// bool status = 1;
// string remarks = 2;
// int32  account_id = 3;
// uint64 balance = 4;
func (s *WalletServiceServer) OpenAccount(ctx context.Context, req *pb.OpenRequest) (*pb.OpenResponse, error) {
	fmt.Println("Registering new wallet account...")
	var acct Account
	acct.AccountID = int(req.AccountId)
	acct.FirstName = req.FirstName
	acct.LastName = req.LastName
	acct.Email = req.Email
	acct.MobileNo = req.MobileNo
	resp, created := acct.NewAcount()
	if !created {
		return &pb.OpenResponse{Remarks: resp, Status: false}, nil
	}

	return &pb.OpenResponse{Remarks: resp, Status: true}, nil
}
func (s *WalletServiceServer) SendTransaction(ctx context.Context, req *pb.TransactionRequest) (*pb.TransactionResponse, error) {
	fmt.Println("Send Wallet Transaction...")
	var trans Transaction
	trans.AccountID = int(req.AccountId)
	trans.Type = req.Type
	trans.Description = req.Description
	trans.Amount = float64(req.Amount)
	trans.SenderID = int(req.SenderId)
	resp, created := trans.Add()
	if !created {
		return &pb.TransactionResponse{Remarks: resp.Remarks, Status: resp.Status}, nil
	}

	return &pb.TransactionResponse{
		Id:              int32(resp.ID),
		AccountId:       int32(resp.AccountID),
		RefNo:           int32(resp.RefNo),
		Description:     resp.Description,
		Amount:          float32(req.Amount),
		Type:            resp.Type,
		Status:          resp.Status,
		Remarks:         resp.Remarks,
		CreatedAt:       resp.CreatedAt.Format("2006-01-02 15:04:05"),
		PreviousBalance: uint64(resp.PreviousBalance),
		CurrentBalance:  uint64(resp.CurrentBalance),
		SenderId:        int32(resp.SenderID),
	}, nil
}

func main() {
	OpenDB()
	fmt.Println("Initializing wallet module...")
	listener, err := net.Listen("tcp", PORT)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Create a new gRPC server
	grpcServer := grpc.NewServer()
	pb.RegisterWalletServiceServer(grpcServer, &WalletServiceServer{})

	log.Println("gRPC server running on port" + PORT)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
	//OpenDB()
}

func OpenDB() {
	fmt.Println("Initializing database wallet databse...")
	os.Setenv("TZ", "Asia/Manila")
	DBM = pg.Connect(&pg.Options{
		Addr:     "localhost:5432",
		User:     "postgres",
		Password: "postgres",
		Database: "hyperball_race_wallet_db",
	})
	var n int
	_, err := DBM.QueryOne(pg.Scan(&n), "SELECT 1")
	if err != nil {
		panic(err)
	} else {
		fmt.Println("Wallet database connected!")
	}
}
