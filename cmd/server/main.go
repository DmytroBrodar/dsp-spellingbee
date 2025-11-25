package main

import (
	"context"
	"fmt"
	"log"
	"net"

	gamepb "spellingbee/api/spellingbee/v1"
	"spellingbee/internal/dictionary"
	"spellingbee/internal/dictionary/decorate"
	"spellingbee/internal/game"
	"spellingbee/internal/manager"

	"google.golang.org/grpc"
)

// our gRPC server type, implements the methods from proto
type server struct {
	gamepb.UnimplementedSpellingBeeServer

	// keep dictionary here but decorator and proxied in main(), this shows dependency injection
	dict dictionary.Dictionary
}

// StartGame creates a new game using Factory, stores it in Singleton manager
// and return the letters to the client
func (s *server) StartGame(ctx context.Context, _ *gamepb.StartRequest) (*gamepb.StartResponse, error) {
	mgr := manager.Get()                             // singleton
	mgr.Game = game.NewGameFromFile("pangrams.json") // factory
	if mgr.Game == nil || len(mgr.Game.Letters) == 0 {
		return nil, fmt.Errorf("failed to create game from file")
	}

	// convert []rune to []string for proto
	letters := make([]string, len(mgr.Game.Letters))
	for i, r := range mgr.Game.Letters {
		letters[i] = string(r)
	}

	return &gamepb.StartResponse{
		Letters: letters,
		Center:  string(mgr.Game.Center),
	}, nil

}

// Submitword checks rules + dictionary and returns score
func (s *server) SubmitWord(ctx context.Context, req *gamepb.WordRequest) (*gamepb.WordResponse, error) {
	mgr := manager.Get()
	if mgr.Game == nil {
		return &gamepb.WordResponse{Valid: false, Message: "Game not started"}, nil
	}

	//local rules
	ok, msg := mgr.Game.IsValidWord(req.GetWord())
	if !ok {
		return &gamepb.WordResponse{Valid: false, Message: msg, Total: int32(mgr.Game.Score)}, nil
	}

	// dictionary check using decorator(logging) and proxy(caching)
	if !s.dict.IsValid(req.GetWord()) {
		return &gamepb.WordResponse{Valid: false, Message: "Not a dictionary word", Total: int32(mgr.Game.Score)}, nil
	}

	// scoring
	points := mgr.Game.ScoreWord(req.GetWord())

	return &gamepb.WordResponse{
		Valid:   true,
		Message: "Valid word",
		Score:   int32(points),
		Total:   int32(mgr.Game.Score),
	}, nil

}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}

	// build dictionary pipelines
	base := dictionary.NewDictionary("words_dictionary.json")
	logged := decorate.WrapLogging(base)  // decoration
	cached := dictionary.NewProxy(logged) // proxy

	s := grpc.NewServer()
	gamepb.RegisterSpellingBeeServer(s, &server{dict: cached})
	log.Println("Server running on :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
