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
	"spellingbee/internal/mq"
	"spellingbee/internal/stats"

	"google.golang.org/grpc"
)

// our gRPC server type, implements the methods from proto
type server struct {
	gamepb.UnimplementedSpellingBeeServer
	dict      dictionary.Dictionary
	stats     *stats.Stats
	statsFile string
	pub       *mq.Publisher
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

// SubmitWord checks rules + dictionary and returns score
func (s *server) SubmitWord(ctx context.Context, req *gamepb.WordRequest) (*gamepb.WordResponse, error) {
	mgr := manager.Get()
	if mgr.Game == nil {
		return &gamepb.WordResponse{Valid: false, Message: "Game not started"}, nil
	}

	word := req.GetWord()

	//local rules
	ok, msg := mgr.Game.IsValidWord(req.GetWord())
	if !ok {
		s.stats.Update(false, 0, false, s.statsFile)
		return &gamepb.WordResponse{Valid: false, Message: msg, Total: int32(mgr.Game.Score)}, nil
	}

	// dictionary check using decorator(logging) and proxy(caching)
	if !s.dict.IsValid(req.GetWord()) {
		s.stats.Update(false, 0, false, s.statsFile)
		return &gamepb.WordResponse{Valid: false, Message: "Not a dictionary word", Total: int32(mgr.Game.Score)}, nil
	}

	// scoring
	points := mgr.Game.ScoreWord(req.GetWord())

	// check pangram
	isPangram := game.IsPangram(word, mgr.Game.Letters)

	// update stats as valid
	s.stats.Update(true, points, isPangram, s.statsFile)

	// if pangram found, publish event to rabbitMq in a goroutine
	if isPangram && s.pub != nil {
		go func(w string) {
			if err := s.pub.PublishPangram(w); err != nil {
				log.Println("failed to publish pangram:", err)
			}
		}(word)
	}

	return &gamepb.WordResponse{
		Valid:   true,
		Message: "Valid word",
		Score:   int32(points),
		Total:   int32(mgr.Game.Score),
	}, nil

}

func main() {
	statsFile := "./internal/stats/stats.txt"

	// load stats from file or create new stats
	st := stats.LoadStats(statsFile)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}

	// build dictionary pipelines
	base := dictionary.NewDictionary("words_dictionary.json")
	logged := decorate.WrapLogging(base)  // decoration
	cached := dictionary.NewProxy(logged) // proxy

	rabbitUrl := "amqp://guest:guest@localhost:5672/"
	exchange := "spellingbee.pangrams"

	pub, err := mq.NewPublisher(rabbitUrl, exchange)
	if err != nil {
		log.Println("Couldn't connect to RabbitMQ:", err)
	}

	// Create server instance WITH stats
	srv := &server{
		dict:      cached,
		stats:     st,
		statsFile: statsFile,
		pub:       pub,
	}

	grpcServer := grpc.NewServer()
	gamepb.RegisterSpellingBeeServer(grpcServer, srv)

	log.Println("Server running on :50051")

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
