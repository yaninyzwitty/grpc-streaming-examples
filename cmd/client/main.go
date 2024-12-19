// grpc client

package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/yaninyzwitty/go-grpc-streaming/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Second)
	defer cancel()
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		slog.Error("failed to create a grpc connection: ", "error", err)
		os.Exit(1)
	}

	defer conn.Close()

	client := pb.NewChatServiceClient(conn)

	stream, err := client.UploadFile(ctx)
	if err != nil {
		slog.Error("failed to create a grpc stream: ", "error", err)
		os.Exit(1)
	}

	fileName := "example.txt"
	for i := 0; i < 10; i++ {
		chunk := &pb.FileChunk{
			Name:        fileName,
			Content:     []byte(fmt.Sprintf("This is chunk %d", i)),
			ChunkNumber: int32(i),
		}

		if err := stream.Send(chunk); err != nil {
			slog.Error("failed to send a message to the server: ", "error", err)
			return
		}
		slog.Info("Sent chunk",
			"chunk_index", i,
			"file_name", fileName,
		)
		time.Sleep(time.Millisecond * 500)

	}

	summary, err := stream.CloseAndRecv()
	if err != nil {
		slog.Error("failed to receive a message from the server: ", "error", err)
		os.Exit(1)
	}
	slog.Info("Summary", "summary", summary)

	// client := pb.NewChatServiceClient(conn)
	// stream, err := client.Chat(ctx)
	// if err != nil {
	// 	slog.Error("failed to create a grpc stream: ", "error", err)
	// 	os.Exit(1)
	// }

	// here we send messages to the server, bi directional

	// go func() {
	// 	for i := 0; i < 5; i++ {
	// 		msg := &pb.ChatMessage{
	// 			User:    "Client",
	// 			Message: "Hello from the client!",
	// 		}
	// 		if err := stream.Send(msg); err != nil {
	// 			slog.Error("failed to send a message to the server: ", "error", err)
	// 			return
	// 		}
	// 	}
	// 	stream.CloseSend()
	// }()

	// for {
	// 	res, err := stream.Recv()
	// 	if err == io.EOF {
	// 		break
	// 	}
	// 	if err != nil {
	// 		slog.Error("failed to receive a message from the server: ", "error", err)
	// 		os.Exit(1)
	// 	}

	// 	slog.Info(" message from the server: ", "message", res.Message)
	// }

}
