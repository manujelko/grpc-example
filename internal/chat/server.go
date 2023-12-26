package chat

import (
	"context"
	"fmt"
	"io"
	"time"

	pb "github.com/manujelko/grpc-example/pkg/api/chat"
)

type server struct {
	pb.UnimplementedSimpleChatServer
}

func NewServer() *server {
	return &server{}
}

func (s *server) Ping(ctx context.Context, in *pb.PingRequest) (*pb.PingResponse, error) {
	return &pb.PingResponse{Message: "pong"}, nil
}

func (s *server) UploadMessages(stream pb.SimpleChat_UploadMessagesServer) error {
	var messages []string
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.UploadStatus{Status: "Received " + fmt.Sprint(len(messages)) + " messages"})
		}
		if err != nil {
			return err
		}
		messages = append(messages, msg.Text)
	}
}

func (s *server) NewsTicker(req *pb.TickerRequest, stream pb.SimpleChat_NewsTickerServer) error {
	for i := 0; i < int(req.NumberOfMessages); i++ {
		news := fmt.Sprintf("News %d", i+1)
		if err := stream.Send(&pb.NewsMessage{NewsText: news}); err != nil {
			return err
		}
		time.Sleep(1 * time.Second)
	}
	return nil
}

func (s *server) EchoChat(stream pb.SimpleChat_EchoChatServer) error {
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		if err := stream.Send(&pb.ChatMessage{Text: "You said: " + in.Text}); err != nil {
			return err
		}
	}
}
