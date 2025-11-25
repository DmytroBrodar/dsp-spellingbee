package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	gamepb "spellingbee/api/spellingbee/v1"
)

// this is the simple text client. it connects to the server, gets letters, and then sends words to check
func main() {
	// connect to server
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	c := gamepb.NewSpellingBeeClient(conn)

	//start game
	start, err := c.StartGame(context.Background(), &gamepb.StartRequest{})
	if err != nil {
		panic(err)
	}

	fmt.Println("---------- Spelling Bee ----------")
	fmt.Printf("Letters: %v (center %s)\n", start.GetLetters(), start.GetCenter())
	fmt.Printf("Type words (Ctrl + C to quit).")

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		line, _ := reader.ReadString('\n')
		word := strings.TrimSpace(line)
		if word == "" {
			continue
		}

		resp, err := c.SubmitWord(context.Background(), &gamepb.WordRequest{Word: word})
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}

		if resp.GetValid() {
			fmt.Printf("Word is valid (+%d) | Total: %d\n", resp.GetScore(), resp.GetTotal())
		} else {
			fmt.Printf("Invalid word: %s | TOTAL: %d\n", resp.GetMessage(), resp.GetTotal())
		}

	}
}
