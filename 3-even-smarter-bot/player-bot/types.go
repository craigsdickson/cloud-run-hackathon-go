package main

import (
	"player-bot/shared"
)

type ArenaUpdate struct {
	Links struct {
		Self struct {
			Href string `json:"href"`
		} `json:"self"`
	} `json:"_links"`
	Arena struct {
		Dimensions []int                         `json:"dims"`
		State      map[string]shared.PlayerState `json:"state"`
	} `json:"arena"`
}

type Coordinates struct {
	X int
	Y int
}
