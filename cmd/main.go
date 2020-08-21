package main

import (
	"fmt"

	"github.com/M4cs/gomodio"
)

func main() {
	user := gomodio.NewUser("c1ac9d8941bb8d42bf4f43472f9ec940", "mabridgland@protonmail.com")
	query := map[string]string{
		"name": "Skater XL",
		"q":    "Skater",
	}
	games, _ := gomodio.GetGames(query, user)
	game, _ := gomodio.GetGame(games.Data[0].ID, nil, user)
	fmt.Println(game.Name)
}
