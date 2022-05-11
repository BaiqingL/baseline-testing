package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"

	pb "github.com/BaiqingL/baseline-testing/internal"
)

var m map[string]int = make(map[string]int)

type server struct {
	pb.UnimplementedListenerServer
}

func (s *server) Add(ctx context.Context, in *pb.AddRequest) (*pb.AddReply, error) {
	log.Printf("Received: %v", in.GetKey())
	return &pb.AddReply{Value: 1}, nil
}

func WordCount(s string) map[string]int {
	words := strings.Fields(s)
	for _, word := range words {
		m[word] += 1
	}
	return m
}

func main() {
	fileByte, err := ioutil.ReadFile("wordcount.txt")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Starting word count...")
		// loop 10 times
		for i := 0; i < 10; i++ {
			start := time.Now()
			mapOut := WordCount(string(fileByte))
			duration := time.Since(start)
			fmt.Println("Duration:", duration)
			words := 0
			for _, value := range mapOut {
				words += value
			}
			fmt.Println("Words:", words)
		}
	}

}
