package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
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
	jsonReq, _ := json.Marshal(v)
	log.Printf("received ArenaUpdate: %s", jsonReq)
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
	opponentCoords := determineNextClosestPlayer(myState, board)
	return determineNextMove(myState, opponentCoords)
}

func determineNextMove(myState PlayerState, opponentCoords []int) string {
	commands := []string{"F", "R", "L"}
	rand := rand2.Intn(3)
	return commands[rand]
}

func determineNextClosestPlayer(myState PlayerState, board [][]bool) []int {
	closestCoords := []int{0, 0}
	closestDistance := -1.0
	for x := range board {
		for y := range board[x] {
			if x == myState.X && y == myState.Y { // skip ourselves
				continue
			}
			if board[x][y] { // if there's a player at this location
				currentDistance := calculateDistance(myState, x, y)
				if closestDistance == -1 || currentDistance < closestDistance {
					closestDistance = currentDistance
					closestCoords = []int{x, y}
				}
			}
		}
	}
	return closestCoords
}

func calculateDistance(myState PlayerState, x2 int, y2 int) float64 {
	x1 := myState.X
	y1 := myState.Y
	return math.Sqrt(math.Pow(float64(x2-x1), 2) + math.Pow(float64(y2-y1), 2))
}

// determines if there is a player in our firing line or not
func someoneIsInFrontOfMe(myState PlayerState, board [][]bool) (result bool) {
	myXcoord := myState.X
	myYcoord := myState.Y
	myDirection := myState.Direction
	boardWidth := len(board)
	boardHeight := len(board[0])
	maxThrowLength := 3
	switch myDirection {
	case "N":
		for i := 1; i <= maxThrowLength; i++ {
			if myYcoord-i >= 0 && board[myXcoord][myYcoord-i] { // check we dont go outside north border
				return true
			}
		}
	case "E":
		for i := 1; i <= maxThrowLength; i++ {
			if myXcoord+i < boardWidth && board[myXcoord+i][myYcoord] { // check we dont go outside the east border
				return true
			}
		}
	case "S":
		for i := 1; i <= maxThrowLength; i++ {
			if myYcoord+i < boardHeight && board[myXcoord][myYcoord+i] { // check we dont go outside the south border
				return true
			}
		}
	default: // "W"
		for i := 1; i <= maxThrowLength; i++ {
			if myXcoord-i >= 0 && board[myXcoord-i][myYcoord] { // check we dont go outside west border
				return true
			}
		}
	}
	return false

}
