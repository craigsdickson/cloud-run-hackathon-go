package main

import (
	"encoding/json"
	"even-smarter-bot/board"
	"even-smarter-bot/playerstate"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
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
	resp := play(v)
	fmt.Fprint(w, resp)
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
		// otherwise check if someone is in front of us and within the max throw distance
	} else if board.IsSomeoneInFrontOfMe(myState, 3) {
		log.Printf("I am throwing")
		return "T"
		// otherwise move towards an opponent we want to throw at
	} else {
		log.Printf("there is no one to throw at")
		//return moveTowardsClosestOpponent(myState, board)
		return moveTowardsClosestHighScoringOpponent(myState, board)
	}
}

func extractMyState(input ArenaUpdate) playerstate.PlayerState {
	myId := input.Links.Self.Href
	state := input.Arena.State
	return state[myId]
}

func moveTowardsClosestOpponent(myState playerstate.PlayerState, board board.Board) (response string) {
	opponent := board.FindClosestOpponent(myState)
	log.Printf("closest opponent is at x:%v y:%v", opponent.X, opponent.Y)
	return determineNextMove(myState, opponent)
}

func moveTowardsClosestHighScoringOpponent(myState playerstate.PlayerState, board board.Board) (response string) {
	opponent := board.FindClosestHighScoringOpponent(myState, 0.5)
	log.Printf("closest high scoring opponent is at x:%v y:%v with a score of %v", opponent.X, opponent.Y, opponent.Score)
	return determineNextMove(myState, opponent)
}

func determineNextMove(myState playerstate.PlayerState, opponentState playerstate.PlayerState) (result string) {
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
	log.Printf("direction of opponent is %v, therefore I am going to move %v", directionOfOpponent, result)
	return result
}

func determineDirectionOfOpponent(myState playerstate.PlayerState, opponentState playerstate.PlayerState) (result string) {
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
