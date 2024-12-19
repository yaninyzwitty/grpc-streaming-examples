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
	"net"
	"os"

	"github.com/yaninyzwitty/go-grpc-streaming/pb"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedChatServiceServer
}

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
