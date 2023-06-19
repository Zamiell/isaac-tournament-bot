package main

// The type of element in the array in the "builds.json" file.
type Build struct {
	Name             string            `json:"name"`
	Category         string            `json:"category"`
	Collectibles     []Collectible     `json:"collectibles"`
	BannedCharacters []BannedCharacter `json:"bannedCharacters"`
}

type Collectible struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type BannedCharacter struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Reason string `json:"reason"`
}
