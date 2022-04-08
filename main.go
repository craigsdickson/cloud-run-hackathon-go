package main

import (
	"encoding/json"
	"fmt"
	"log"
	rand2 "math/rand"
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
	board := generateBoard(input)
	myState := extractMyState(input)
	if someoneIsInFrontOfMe(myState, board) {
		log.Println("throwing because someone is in front of me")
		return "T"
	} else {
		return moveTowardsNextClosestPlayer(myState, board)
	}
}

func extractMyState(input ArenaUpdate) PlayerState {
	myId := input.Links.Self.Href
	state := input.Arena.State
	return state[myId]
}

// creates a 2D array of booleans representing locations of players on the board
func generateBoard(input ArenaUpdate) [][]bool {
	// generate the board data structure
	width := input.Arena.Dimensions[0]
	height := input.Arena.Dimensions[1]
	var board = make([][]bool, width)
	for i := range board {
		board[i] = make([]bool, height)
	}
	// now populate board with locations of players
	state := input.Arena.State
	for _, v := range state {
		vX := v.X
		vY := v.Y
		board[vX][vY] = true
	}
	log.Printf("board is: %v", board)
	return board
}

func moveTowardsNextClosestPlayer(myState PlayerState, board [][]bool) (response string) {
	commands := []string{"F", "R", "L", "T"}
	rand := rand2.Intn(4)
	return commands[rand]
}

// determines if there is a player in firing line or not
func someoneIsInFrontOfMe(myState PlayerState, board [][]bool) (result bool) {
	myX := myState.X
	myY := myState.Y
	myDirection := myState.Direction
	boardWidth := len(board)
	boardHeight := len(board[0])
	maxThrowLength := 3
	switch myDirection {
	case "N":
		for i := 1; i <= maxThrowLength; i++ {
			if myY-i >= 0 && board[myX][myY-i] { // check we dont go outside north border
				return true
			}
		}
	case "E":
		for i := 1; i <= maxThrowLength; i++ {
			if myX+i < boardWidth && board[myX+i][myY] { // check we dont go outside the east border
				return true
			}
		}
	case "S":
		for i := 1; i <= maxThrowLength; i++ {
			if myY+i <= boardHeight && board[myX][myY+i] { // check we dont go outside the south border
				return true
			}
		}
	default: // "W"
		for i := 1; i <= maxThrowLength; i++ {
			if myX-i >= 0 && board[myX-i][myY] { // check we dont go outside west border
				return true
			}
		}
	}
	return false

}
