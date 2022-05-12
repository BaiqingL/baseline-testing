package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"strings"
	"time"

	"github.com/BaiqingL/baseline-testing/internal/upstream"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var m map[string]int = make(map[string]int)

type server struct {
	upstream.UnimplementedListenerServer
}

func (s *server) Add(ctx context.Context, in *upstream.AddRequest) (*upstream.AddResponse, error) {
	m[in.GetKey()] += 1
	return &upstream.AddResponse{Value: 1}, nil
}
func startServer() {
	// make a channel
	s := grpc.NewServer()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 8080))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	upstream.RegisterListenerServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
	fmt.Println("server started")
}

func main() {
	go startServer()

	// client
	wordBank, err := ioutil.ReadFile("wordcount.txt")
	words := strings.Fields(string(wordBank))
	if err != nil {
		log.Fatalf("failed to read file: %v", err)
	}
	conn, err := grpc.Dial("127.0.0.1:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()
	client := upstream.NewListenerClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	fmt.Println("client started")
	start := time.Now()
	idx := 0
	for _, word := range words {
		_, err := client.Add(ctx, &upstream.AddRequest{Key: word, Value: 1})
		idx++
		if err != nil {
			log.Fatalf("could not send: %v", err)
		}
		if (idx % 1000) == 0 {
			fmt.Println(idx)
		}
	}
	duration := time.Since(start)
	fmt.Println("Duration:", duration)
}
