package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"

	"internal/protoCustom"
)

type server struct {
	protoCustom.UnimplementedListenerServer
}

func (s *server) SayHello(ctx context.Context, in *protoCustom.AddRequest) (*protoCustom.AddReply, error) {
	log.Printf("Received: %v", in.GetName())
	return &protoCustom.AddReply{Value: 1}, nil
}

func WordCount(s string) map[string]int {
	words := strings.Fields(s)
	m := make(map[string]int)
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
