package integration

import (
	"context"
	"fmt"
	"io"
	"testing"

	pb "github.com/manujelko/grpc-example/pkg/api/chat"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestPing(t *testing.T) {
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to dial server: %v", err)
	}
	defer conn.Close()

	client := pb.NewSimpleChatClient(conn)

	resp, err := client.Ping(context.Background(), &pb.PingRequest{Message: "ping"})
	if err != nil {
		t.Fatalf("Ping failed: %v", err)
	}

	if resp.GetMessage() != "pong" {
		t.Errorf("Expected 'pong', got '%s'", resp.GetMessage())
	}
}

func TestUploadMessages(t *testing.T) {
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to dial server: %v", err)
	}
	defer conn.Close()

	client := pb.NewSimpleChatClient(conn)

	stream, err := client.UploadMessages(context.Background())
	if err != nil {
		t.Fatalf("Error creating stream: %v", err)
	}

	messages := []string{"Hello", "World", "Test"}
	for _, msg := range messages {
		if err := stream.Send(&pb.Message{Text: msg}); err != nil {
			t.Fatalf("Failed to send message: %v", err)
		}
	}

	resp, err := stream.CloseAndRecv()
	if err != nil {
		t.Fatalf("Error receiving response: %v", err)
	}

	expectedStatus := "Received 3 messages"
	if resp.GetStatus() != expectedStatus {
		t.Errorf("Expected '%s', got '%s'", expectedStatus, resp.GetStatus())
	}
}

func TestNewsTicker(t *testing.T) {
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to dial server: %v", err)
	}
	defer conn.Close()

	client := pb.NewSimpleChatClient(conn)

	stream, err := client.NewsTicker(context.Background(), &pb.TickerRequest{NumberOfMessages: 5})
	if err != nil {
		t.Fatalf("Error creating stream: %v", err)
	}
	for i := 1; i <= 5; i++ {
		msg, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("Failed to receive message: %v", err)
		}

		expectedText := fmt.Sprintf("News %d", i)
		if msg.GetNewsText() != expectedText {
			t.Errorf("Expected '%s', got '%s'", expectedText, msg.GetNewsText())
		}
	}
}

func TestEchoChat(t *testing.T) {
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to dial server: %v", err)
	}
	defer conn.Close()

	client := pb.NewSimpleChatClient(conn)

	stream, err := client.EchoChat(context.Background())
	if err != nil {
		t.Fatalf("Error creating stream: %v", err)
	}

	messages := []string{"Hi", "How are you?", "Test"}
	for _, msg := range messages {
		if err := stream.Send(&pb.ChatMessage{Text: msg}); err != nil {
			t.Fatalf("Failed to send message: %v", err)
		}

		reply, err := stream.Recv()
		if err != nil {
			t.Fatalf("Failed to receive message: %v", err)
		}

		expectedText := "You said: " + msg
		if reply.GetText() != expectedText {
			t.Errorf("Expected '%s', got '%s'", expectedText, reply.GetText())
		}
	}

	stream.CloseSend()
}
