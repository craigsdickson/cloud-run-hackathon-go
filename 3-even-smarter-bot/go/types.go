package main

import (
	"even-smarter-bot/playerstate"
)

type ArenaUpdate struct {
	Links struct {
		Self struct {
			Href string `json:"href"`
		} `json:"self"`
	} `json:"_links"`
	Arena struct {
		Dimensions []int                              `json:"dims"`
		State      map[string]playerstate.PlayerState `json:"state"`
	} `json:"arena"`
}

type Coordinates struct {
	X int
	Y int
}
