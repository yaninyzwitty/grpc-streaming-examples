// grpc streaming, sending a stream of messages
// this is diff from single request to single response fashion such in unary rpc
// grpc includes, 1, client streaming 2, server streaming 3, bidirectional streaming
// 1, client streaming, client sends a stream of messages to server, server receives a stream of messages
// 2, server streaming, server sends a stream of messages to client, client receives a stream of messages
// 3, bidirectional streaming, client sends a stream of messages to server, server sends a stream of messages to client, client receives a stream of messages
// 4, unary rpc, client sends a single request to server, server sends a single response to client, client receives a single response

package main

import (
	"fmt"
	"io"
	"log"
	"log/slog"
	"math/rand"
	"net"
	"os"
	"time"

	"github.com/yaninyzwitty/go-grpc-streaming/pb"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedChatServiceServer
}

// bi-directional
func (s *server) Chat(stream pb.ChatService_ChatServer) error {
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return fmt.Errorf("error receiving messages: %w", err)
		}

		slog.Info("Received message", "message", msg.Message)
		slog.Info("Receiving user", "user", msg.User)

		// Send a response back to the client
		err = stream.Send(&pb.ChatMessage{
			User:    "Witty",
			Message: fmt.Sprintf("Hello %s, you said: %s", msg.User, msg.Message),
		})

		if err != nil {
			log.Fatalf("Error sending message: %v", err)
		}

	}
}

// client streaming

func (s *server) UploadFile(stream pb.ChatService_UploadFileServer) error {
	totalChunks := 0
	fileName := ""
	for {
		chunk, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.UploadSummary{
				Message:     fmt.Sprintf("File %s uploaded successfully with %d chunks", fileName, totalChunks),
				TotalChunks: int32(totalChunks),
			})

		}
		if err != nil {
			slog.Error("Error receiving chunk", "error", err)
			return err
		}

		// Process the chunk

		fileName = chunk.Name
		totalChunks++
		slog.Info("Received chunk", "chunk", "name", fmt.Sprintf("%d", chunk.ChunkNumber), fileName)
	}
}

// Simulate a server streaming RPC for real-time stock prices.

func (s *server) GetStockPrices(req *pb.StockRequest, stream pb.ChatService_GetStockPricesServer) error {

	symbol := req.Symbol
	slog.Info("Received request for stock prices", "symbol", symbol)

	// Simulate sending a stream of stock prices

	for i := 0; i < 10; i++ {
		price := rand.Float32() * 600
		timestamp := time.Now().Unix()

		resp := &pb.StockResponse{
			Symbol:    symbol,
			Price:     price,
			Timestamp: timestamp,
		}

		slog.Info("Sending stock price", "symbol", symbol, "price", price, "timestamp", timestamp)
		if err := stream.Send(resp); err != nil {
			slog.Error("Error sending stock price", "error", err)
			return err
		}
		time.Sleep(time.Second)

	}
	return nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		slog.Error("failed to listen", "error", err)
		os.Exit(1)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterChatServiceServer(grpcServer, &server{})
	slog.Info("Server listening on port 50051...")
	grpcServer.Serve(lis)

}
