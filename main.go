package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func main() {
	response, err := http.Get("https://api.chess.com/pub/player/ferenco/games/2023/07")
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		log.Fatalf("Unexpected response status code: %s", response.Status)
	}

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	var monthlyChessGames MonthlyChessGames
	json.Unmarshal(responseData, &monthlyChessGames)

	for _, game := range monthlyChessGames.Games {
		fmt.Printf("White: %s\n", game.White.Username)
		fmt.Printf("Black: %s\n", game.Black.Username)
		if game.White.Result == "win" {
			fmt.Printf("Winner: %s with a rating of: %d\n", game.White.Username, game.White.Rating)
		} else {
			fmt.Printf("Winner: %s with a rating of: %d\n", game.Black.Username, game.Black.Rating)
		}
		fmt.Println()
	}
}

// A MonthlyChessGames contains each game a user played in a given Month.
type MonthlyChessGames struct {
	Games []struct {
		URL          string `json:"url"`
		Pgn          string `json:"pgn"`
		TimeControl  string `json:"time_control"`
		EndTime      int    `json:"end_time"`
		Rated        bool   `json:"rated"`
		Tcn          string `json:"tcn"`
		UUID         string `json:"uuid"`
		InitialSetup string `json:"initial_setup"`
		Fen          string `json:"fen"`
		TimeClass    string `json:"time_class"`
		Rules        string `json:"rules"`
		White        struct {
			Rating   int    `json:"rating"`
			Result   string `json:"result"`
			ID       string `json:"@id"`
			Username string `json:"username"`
			UUID     string `json:"uuid"`
		} `json:"white"`
		Black struct {
			Rating   int    `json:"rating"`
			Result   string `json:"result"`
			ID       string `json:"@id"`
			Username string `json:"username"`
			UUID     string `json:"uuid"`
		} `json:"black"`
		Accuracies struct {
			White float64 `json:"white"`
			Black float64 `json:"black"`
		} `json:"accuracies,omitempty"`
	} `json:"games"`
}
