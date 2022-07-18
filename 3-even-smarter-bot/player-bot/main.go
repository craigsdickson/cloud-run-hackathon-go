package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"player-bot/board"
	"player-bot/shared"

	"cloud.google.com/go/compute/metadata"
	"cloud.google.com/go/pubsub"

	"github.com/gomodule/redigo/redis"
)

var redisPool *redis.Pool

var HIGH_SCORING_PERCENTILE = 0.5
var MAX_THROW_DISTANCE = 3

func main() {
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	redisAddr := fmt.Sprintf("%s:%s", redisHost, redisPort)

	const maxConnections = 10
	redisPool = &redis.Pool{
		MaxIdle: maxConnections,
		Dial:    func() (redis.Conn, error) { return redis.Dial("tcp", redisAddr) },
	}

	port := "8080"
	if v := os.Getenv("PORT"); v != "" {
		port = v
	}
	http.HandleFunc("/", handler)

	log.Printf("starting server on port :%v", port)
	err := http.ListenAndServe(":"+port, nil)
	log.Fatalf("http listen error: %v", err)
}

func handler(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodGet {
		fmt.Fprint(w, "Let the battle begin!")
		return
	}

	var v ArenaUpdate
	defer req.Body.Close()
	d := json.NewDecoder(req.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(&v); err != nil {
		log.Printf("WARN: failed to decode ArenaUpdate in response body: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	go postArenaUpdateEvent(v) // call this asynchonously
	resp := play(v)
	fmt.Fprint(w, resp)
}

func postArenaUpdateEvent(input ArenaUpdate) {
	ctx := context.Background()
	metadataClient := metadata.NewClient(nil)
	projectId, err := metadataClient.ProjectID()
	if err != nil {
		log.Fatalf("metadataClient.ProjectID: %v", err)
	}
	pubsubClient, err := pubsub.NewClient(ctx, projectId)
	if err != nil {
		log.Fatalf("pubsub.NewClient: %v", err)
	}
	defer pubsubClient.Close()
	topic := pubsubClient.Topic(os.Getenv("ARENA_UPDATES_PUBSUB_TOPIC_NAME"))
	message, err := json.Marshal(input)
	if err != nil {
		log.Fatalf("json.Marshal fatal error: %v", err)
	}
	publishResult := topic.Publish(ctx, &pubsub.Message{Data: message})
	id, err := publishResult.Get(ctx)
	if err != nil {
		log.Fatalf("topic.Publish fatal error: %v", err)
	}
	log.Printf("published message with id: %v", id)
	topic.Stop()
}

func play(input ArenaUpdate) (response string) {
	log.Printf("IN: %#v", input)
	board := board.New(input.Arena.Dimensions[0], input.Arena.Dimensions[1], input.Arena.State)
	myState := extractMyState(input)
	log.Printf("i am at x:%v y:%v and I am facing %v", myState.X, myState.Y, myState.Direction)
	// if we are the only player, just spin on the spot
	if board.NumberOfPlayers == 1 {
		log.Printf("there are no other players on the board")
		return "R"
	}
	// check to see if there is a leaderboard available, otherwsie just look for closest player
	leaderboard := getLeaderboard()
	if leaderboard != nil {
		if board.IsThereAHighScoringOpponentInFrontOfMe(myState, MAX_THROW_DISTANCE, leaderboard, HIGH_SCORING_PERCENTILE) {
			log.Printf("there is a highscoring opponent in front of me, so I am going to throw")
			return "T"
		} else {
			log.Printf("there are no highscoring opponents to throw at")
			return moveTowardsClosestHighScoringOpponent(myState, board, leaderboard)
		}
	} else {
		if board.IsThereAnOpponentInFrontOfMe(myState, MAX_THROW_DISTANCE) {
			log.Printf("there is an opponent in front of me, so I am throwing")
			return "T"
		} else {
			log.Printf("there are no opponents to throw at")
			return moveTowardsClosestOpponent(myState, board)
		}
	}
}

func getLeaderboard() []shared.PlayerState {
	conn := redisPool.Get()
	defer conn.Close()
	leaderboardAsString, err := redis.String(conn.Do("GET", "leaderboard"))
	if err != nil {
		log.Printf("error reading leaderboard from redis: %v", err)
		return nil
	}
	var leaderboard []shared.PlayerState
	err = json.Unmarshal([]byte(leaderboardAsString), &leaderboard)
	if err != nil {
		log.Printf("error unmarshalling leaderboard: %v", err)
		return nil
	}
	log.Printf("leaderboard is: %v", leaderboard)
	return leaderboard
}

func extractMyState(input ArenaUpdate) shared.PlayerState {
	myId := input.Links.Self.Href
	state := input.Arena.State
	return state[myId]
}

func moveTowardsClosestOpponent(myState shared.PlayerState, board board.Board) (response string) {
	opponent := board.FindClosestOpponent(myState)
	log.Printf("closest opponent is at x:%v y:%v", opponent.X, opponent.Y)
	return determineNextMove(myState, opponent)
}

func moveTowardsClosestHighScoringOpponent(myState shared.PlayerState, board board.Board, leaderboard []shared.PlayerState) (response string) {
	opponent := board.FindClosestHighScoringOpponent(myState, leaderboard, HIGH_SCORING_PERCENTILE)
	log.Printf("closest high scoring opponent is at x:%v y:%v with a score of %v", opponent.X, opponent.Y, opponent.Score)
	return determineNextMove(myState, opponent)
}

func determineNextMove(myState shared.PlayerState, opponentState shared.PlayerState) (result string) {
	directionImFacing := myState.Direction
	directionOfOpponent := determineDirectionOfOpponent(myState, opponentState)
	switch directionImFacing {
	case "N":
		switch directionOfOpponent {
		case "N":
			fallthrough
		case "NE":
			fallthrough
		case "NW":
			result = "F"
		case "E":
			fallthrough
		case "SE":
			fallthrough
		case "S":
			result = "R"
		case "SW":
			fallthrough
		default: // "W":
			result = "L"
		}
	case "E":
		switch directionOfOpponent {
		case "N":
			fallthrough
		case "NW":
			result = "L"
		case "NE":
			fallthrough
		case "E":
			fallthrough
		case "SE":
			result = "F"
		case "S":
			fallthrough
		case "SW":
			fallthrough
		default: // "W":
			result = "R"
		}
	case "S":
		switch directionOfOpponent {
		case "N":
			fallthrough
		case "W":
			fallthrough
		case "NW":
			result = "R"
		case "NE":
			fallthrough
		case "E":
			result = "L"
		case "SE":
			fallthrough
		case "S":
			fallthrough
		default: // "SW":
			result = "F"
		}
	default: //W
		switch directionOfOpponent {
		case "N":
			fallthrough
		case "NE":
			fallthrough
		case "E":
			result = "R"
		case "SE":
			fallthrough
		case "S":
			result = "L"
		case "SW":
			fallthrough
		case "W":
			fallthrough
		default: // "NW":
			result = "F"
		}
	}
	log.Printf("direction I am facing is %v, direction of opponent is %v, therefore I am going to move %v", directionImFacing, directionOfOpponent, result)
	return result
}

func determineDirectionOfOpponent(myState shared.PlayerState, opponentState shared.PlayerState) (result string) {
	myXcoord := myState.X
	myYcoord := myState.Y
	opponentXcoord := opponentState.X
	opponentYcoord := opponentState.Y
	if myXcoord == opponentXcoord {
		if myYcoord > opponentYcoord {
			result = "N"
		} else {
			result = "S"
		}
	} else if myYcoord == opponentYcoord {
		if myXcoord > opponentXcoord {
			result = "W"
		} else {
			result = "E"
		}
	} else if myYcoord > opponentYcoord {
		if myXcoord > opponentXcoord {
			result = "NW"
		} else {
			result = "NE"
		}
	} else { // myYcoord < opponentYcoord
		if myXcoord > opponentXcoord {
			result = "SW"
		} else {
			result = "SE"
		}
	}
	log.Printf("i am at x:%v y:%v and opponent is at x:%v y:%v, so their direction from me is %v", myXcoord, myYcoord, opponentXcoord, opponentYcoord, result)
	return result
}
