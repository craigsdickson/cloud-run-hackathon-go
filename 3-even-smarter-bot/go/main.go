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

	log.Printf("starting server on port :%s", port)
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
	if board.IsSomeoneInFrontOfMe(myState, 3) {
		log.Println("throwing because someone is in front of me")
		return "T"
	} else {
		return moveTowardsClosestOpponent(myState, board)
	}
}

func extractMyState(input ArenaUpdate) playerstate.PlayerState {
	myId := input.Links.Self.Href
	state := input.Arena.State
	return state[myId]
}

func moveTowardsClosestOpponent(myState playerstate.PlayerState, board board.Board) (response string) {
	opponent := board.FindClosestOpponent(myState)
	return determineNextMove(myState, opponent)
}

func moveTowardsClosestHighScoringOpponent(myState playerstate.PlayerState, board board.Board) (response string) {
	opponent := board.FindClosestHighScoringOpponent(myState, 0.5)
	return determineNextMove(myState, opponent)
}

func determineNextMove(myState playerstate.PlayerState, opponentState playerstate.PlayerState) string {
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
			return "F"
		case "E":
			fallthrough
		case "SE":
			fallthrough
		case "S":
			return "R"
		case "SW":
			fallthrough
		default: // "W":
			return "L"
		}
	case "E":
		switch directionOfOpponent {
		case "N":
			fallthrough
		case "NW":
			return "L"
		case "NE":
			fallthrough
		case "E":
			fallthrough
		case "SE":
			return "F"
		case "S":
			fallthrough
		case "SW":
			fallthrough
		default: // "W":
			return "R"
		}
	case "S":
		switch directionOfOpponent {
		case "N":
			fallthrough
		case "W":
			fallthrough
		case "NW":
			return "R"
		case "NE":
			fallthrough
		case "E":
			return "L"
		case "SE":
			fallthrough
		case "S":
			fallthrough
		default: // "SW":
			return "F"
		}
	default: //W
		switch directionOfOpponent {
		case "N":
			fallthrough
		case "NE":
			fallthrough
		case "E":
			return "R"
		case "SE":
			fallthrough
		case "S":
			return "L"
		case "SW":
			fallthrough
		case "W":
			fallthrough
		default: // "NW":
			return "F"
		}
	}
}

func determineDirectionOfOpponent(myState playerstate.PlayerState, opponentState playerstate.PlayerState) string {
	myXcoord := myState.X
	myYcoord := myState.Y
	result := ""
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
	return result
}
