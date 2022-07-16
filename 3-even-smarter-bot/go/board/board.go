package board

import (
	"even-smarter-bot/playerstate"
	"log"
	"math"
	"sort"
)

type Board struct {
	Squares         [][]*playerstate.PlayerState
	Width           int
	Height          int
	NumberOfPlayers int
	Leaderboard     []*playerstate.PlayerState
}

func New(width int, height int, players map[string]playerstate.PlayerState) Board {
	board := Board{}
	board.Width = width
	board.Height = height
	board.NumberOfPlayers = len(players)
	board.Leaderboard = make([]*playerstate.PlayerState, board.NumberOfPlayers)
	board.Squares = make([][]*playerstate.PlayerState, width)
	for i := range board.Squares {
		board.Squares[i] = make([]*playerstate.PlayerState, height)
	}
	// now populate squares and leaderboard with players
	var playerIndex = 0
	for _, v := range players {
		vX := v.X
		vY := v.Y
		board.Squares[vX][vY] = &v
		board.Leaderboard[playerIndex] = &v
		playerIndex++
	}
	// now sort the leaderboard
	sort.Slice(board.Leaderboard, func(i, j int) bool {
		return board.Leaderboard[i].Score < board.Leaderboard[j].Score
	})
	log.Printf("board is: %v", board)
	return board
}

func (board Board) IsSquareOccupied(x int, y int) bool {
	return board.Squares[x][y] != nil
}

// determines if there is a player in front of provided player, within the max distance
func (board Board) IsSomeoneInFrontOfMe(myState playerstate.PlayerState, maxDistance int) (result bool) {
	myXcoord := myState.X
	myYcoord := myState.Y
	myDirection := myState.Direction
	switch myDirection {
	case "N":
		for i := 1; i <= maxDistance; i++ {
			if myYcoord-i >= 0 && board.IsSquareOccupied(myXcoord, myYcoord-i) { // check we dont go outside north border
				return true
			}
		}
	case "E":
		for i := 1; i <= maxDistance; i++ {
			if myXcoord+i < board.Width && board.IsSquareOccupied(myXcoord+i, myYcoord) { // check we dont go outside the east border
				return true
			}
		}
	case "S":
		for i := 1; i <= maxDistance; i++ {
			if myYcoord+i < board.Height && board.IsSquareOccupied(myXcoord, myYcoord+i) { // check we dont go outside the south border
				return true
			}
		}
	default: // "W"
		for i := 1; i <= maxDistance; i++ {
			if myXcoord-i >= 0 && board.IsSquareOccupied(myXcoord-i, myYcoord) { // check we dont go outside west border
				return true
			}
		}
	}
	return false

}

// TODO: optimise to search concentrically out from the players location instead of scanning whole board
func (board Board) FindClosestOpponent(myState playerstate.PlayerState) playerstate.PlayerState {
	closestOpponent := playerstate.PlayerState{}
	closestDistance := -1.0
	for x := range board.Squares {
		for y := range board.Squares[x] {
			if x == myState.X && y == myState.Y { // skip ourselves
				continue
			}
			if board.IsSquareOccupied(x, y) {
				currentDistance := calculateDistance(myState.X, myState.Y, x, y)
				if closestDistance == -1 || currentDistance < closestDistance {
					closestDistance = currentDistance
					closestOpponent = *board.Squares[x][y]
				}
			}
		}
	}
	return closestOpponent
}

/**
 * The percentile controls what makes a player a "high scorer" - a value of 0.1 means only the top 10% of scoring players count,
 * a value of 0.5 means the top 50% of players count etc.
 */
func (board Board) FindClosestHighScoringOpponent(myState playerstate.PlayerState, percentile float64) playerstate.PlayerState {
	closestHighScoringOpponent := playerstate.PlayerState{}
	closestDistance := math.MaxFloat64 // technically this means this method could fail with an incredibly huge board
	var maxIndex int = int(math.Round(float64(board.NumberOfPlayers) * percentile))
	foundAnOpponent := false
	for i := 0; i < maxIndex; i++ {
		opponent := *board.Leaderboard[i]
		if opponent != myState {
			currentDistance := calculateDistance(myState.X, myState.Y, opponent.X, opponent.Y)
			// we save the firstdss opponent, after that we only save the opponent if they are actually closer
			if !foundAnOpponent || currentDistance < closestDistance {
				closestDistance = currentDistance
				closestHighScoringOpponent = opponent
				foundAnOpponent = true
			}
		}
	}

	// for x := range board.Squares {
	// 	for y := range board.Squares[x] {
	// 		if x == myState.X && y == myState.Y { // skip ourselves
	// 			continue
	// 		}
	// 		if board.IsSquareOccupied(x, y) {
	// 			currentDistance := calculateDistance(myState.X, myState.Y, x, y)
	// 			if closestDistance == -1 || currentDistance < closestDistance {
	// 				closestDistance = currentDistance
	// 				closestHighScoringOpponent = *board.Squares[x][y]
	// 			}
	// 		}
	// 	}
	// }
	return closestHighScoringOpponent
}

func calculateDistance(x1 int, y1 int, x2 int, y2 int) float64 {
	return math.Sqrt(math.Pow(float64(x2-x1), 2) + math.Pow(float64(y2-y1), 2))
}
