package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	// rand2 "math/rand"
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
	directionImFacing := myState.Direction
	directionOfOpponent := determineDirectionOfOpponent(myState, opponentCoords)
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

	// switch directionOfOpponent {
	// case "N":
	// 	switch directionImFacing {
	// 	case "N":
	// 		return "F"
	// 	case "E":
	// 		return "R"
	// 	case "S":
	// 		return "R"
	// 	default: //W
	// 		return "L"
	// 	}
	// case "E":
	// 	switch directionImFacing {
	// 	case "N":
	// 		return "R"
	// 	case "E":
	// 		return "F"
	// 	case "S":
	// 		return "R"
	// 	default: //W
	// 		return "L"
	// 	}
	// case "S":
	// 	switch directionImFacing {
	// 	case "N":
	// 		return "R"
	// 	case "E":
	// 		return "R"
	// 	case "S":
	// 		return "F"
	// 	default: //W
	// 		return "L"
	// 	}
	// default: //W
	// 	switch directionImFacing {
	// 	case "N":
	// 		return "L"
	// 	case "E":
	// 		return "R"
	// 	case "S":
	// 		return "R"
	// 	default: //W
	// 		return "F"
	// 	}
	// }
	// commands := []string{"F", "R", "L"}
	// rand := rand2.Intn(3)
	// return commands[rand]
}

func determineDirectionOfOpponent(myState PlayerState, opponentCoords []int) string {
	myXcoord := myState.X
	myYcoord := myState.Y
	result := ""
	opponentXcoord := opponentCoords[0]
	opponentYcoord := opponentCoords[1]
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
