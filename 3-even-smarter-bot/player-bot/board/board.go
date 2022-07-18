package board

import (
	"log"
	"math"
	"player-bot/shared"
)

type Board struct {
	Squares         [][]*shared.PlayerState
	Width           int
	Height          int
	NumberOfPlayers int
	// Leaderboard     []*shared.PlayerState
}

func New(width int, height int, players map[string]shared.PlayerState) Board {
	board := Board{}
	board.Width = width
	board.Height = height
	board.NumberOfPlayers = len(players)
	// board.Leaderboard = make([]*shared.PlayerState, board.NumberOfPlayers)
	board.Squares = make([][]*shared.PlayerState, width)
	for i := range board.Squares {
		board.Squares[i] = make([]*shared.PlayerState, height)
	}

	// now populate squares and leaderboard with players
	var playerIndex = 0
	for _, v := range players {
		vX := v.X
		vY := v.Y
		board.Squares[vX][vY] = &v
		// board.Leaderboard[playerIndex] = &v
		playerIndex++
	}
	// now sort the leaderboard
	// sort.Slice(board.Leaderboard, func(i, j int) bool {
	// 	return board.Leaderboard[i].Score < board.Leaderboard[j].Score
	// })
	log.Printf("board is: %v", board)
	return board
}

func (board Board) IsSquareOccupied(x int, y int) bool {
	return board.Squares[x][y] != nil
}

func (board Board) IsSquareOccupiedByTargetOpponents(x int, y int, targetOpponents []shared.PlayerState) bool {
	if board.IsSquareOccupied(x, y) {
		for _, targetOpponent := range targetOpponents {
			if *board.Squares[x][y] == targetOpponent {
				return true
			}
		}
	}
	return false
}

// determines if there is an opponent in front of provided player, within the max distance
func (board Board) IsThereAnOpponentInFrontOfMe(myState shared.PlayerState, maxDistance int) (result bool) {
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

// determines if there is a high scoring opponent in front of provided player, within the max distance
func (board Board) IsThereAHighScoringOpponentInFrontOfMe(myState shared.PlayerState, maxDistance int, leaderboard []shared.PlayerState, percentile float64) (result bool) {
	myXcoord := myState.X
	myYcoord := myState.Y
	myDirection := myState.Direction
	highScoringOpponents := getHighScoringOpponents(myState, leaderboard, percentile)
	switch myDirection {
	case "N":
		for i := 1; i <= maxDistance; i++ {
			if myYcoord-i >= 0 && board.IsSquareOccupiedByTargetOpponents(myXcoord, myYcoord-i, highScoringOpponents) { // check we dont go outside north border
				return true
			}
		}
	case "E":
		for i := 1; i <= maxDistance; i++ {
			if myXcoord+i < board.Width && board.IsSquareOccupiedByTargetOpponents(myXcoord+i, myYcoord, highScoringOpponents) { // check we dont go outside the east border
				return true
			}
		}
	case "S":
		for i := 1; i <= maxDistance; i++ {
			if myYcoord+i < board.Height && board.IsSquareOccupiedByTargetOpponents(myXcoord, myYcoord+i, highScoringOpponents) { // check we dont go outside the south border
				return true
			}
		}
	default: // "W"
		for i := 1; i <= maxDistance; i++ {
			if myXcoord-i >= 0 && board.IsSquareOccupiedByTargetOpponents(myXcoord-i, myYcoord, highScoringOpponents) { // check we dont go outside west border
				return true
			}
		}
	}
	return false
}

// TODO: optimise to search concentrically out from the player's location instead of scanning whole board
func (board Board) FindClosestOpponent(myState shared.PlayerState) shared.PlayerState {
	closestOpponent := shared.PlayerState{}
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
	log.Printf("returning closest opponent: %v", closestOpponent)
	return closestOpponent
}

/**
 * The percentile controls what makes a player a "high scorer" - a value of 0.1 means only the top 10% of scoring players count,
 * a value of 0.5 means the top 50% of players count etc.
 */
func (board Board) FindClosestHighScoringOpponent(myState shared.PlayerState, leaderboard []shared.PlayerState, percentile float64) shared.PlayerState {
	closestHighScoringOpponent := shared.PlayerState{}
	closestDistance := math.MaxFloat64 // technically this means this method could fail with an incredibly huge board
	highScoringOpponents := getHighScoringOpponents(myState, leaderboard, percentile)
	for i := 0; i < len(highScoringOpponents); i++ {
		opponent := highScoringOpponents[i]
		currentDistance := calculateDistance(myState.X, myState.Y, opponent.X, opponent.Y)
		if currentDistance < closestDistance {
			closestDistance = currentDistance
			closestHighScoringOpponent = opponent
		}
	}
	log.Printf("returning closest highest scoring opponent: %v", closestHighScoringOpponent)
	return closestHighScoringOpponent
}

func calculateDistance(x1 int, y1 int, x2 int, y2 int) float64 {
	return math.Sqrt(math.Pow(float64(x2-x1), 2) + math.Pow(float64(y2-y1), 2))
}

// Need to test for all the edge/corner cases or no leaderboard, current player being only player on the leaderboard, current player being a leader, current player not being a leader
func getHighScoringOpponents(myState shared.PlayerState, leaderboard []shared.PlayerState, percentile float64) (result []shared.PlayerState) {
	log.Printf("determinig high scoring opponents: my score is: %v, leaderboard length is %v, percentile is: %v", myState.Score, len(leaderboard), percentile)
	var maxIndex int = int(math.Round(float64(len(leaderboard)) * percentile))
	for i := 0; i < maxIndex; i++ {
		if leaderboard[i] != myState { // skip ourselves in case we are a high scorer
			result = append(result, leaderboard[i])
		}

	}
	log.Printf("there are %v high scoring opponents: %v", len(result), result)
	return result
}
