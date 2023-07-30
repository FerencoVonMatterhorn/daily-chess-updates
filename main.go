package main

import (
	"encoding/json"
	"fmt"
	"github.com/charmbracelet/log"
	"io"
	"net/http"
	"time"
)

func main() {
	response, err := http.Get("https://api.chess.com/pub/player/ferenco/games/2023/07")
	if err != nil {
		log.Fatal(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(response.Body)

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		log.Fatalf("Unexpected response status code: %s", response.Status)
	}

	b, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatalf("Error decoding API response: %s", err)
	}

	// Create a variable to hold the API response data
	var games Games

	// Decode the API response JSON into the games variable
	err = json.Unmarshal(b, &games)
	if err != nil {
		log.Fatalf("Error decoding API response: %s", err)
	}

	// Print the list of games
	for _, game := range getAllGamesForDate(games.Games, 30, time.July, 2023) {
		printBasicInformation(game)
	}
}

// printBasicInformation prints some basic information about a game (debugging purposes)
func printBasicInformation(game Game) {
	log.Infof("Game played at: %v", getTimeForUnixTimestamp(game.EndTime))
	log.Infof("White: %s", game.White.Username)
	log.Infof("Black: %s", game.Black.Username)
	if game.White.Result == "win" {
		log.Infof("Winner: %s with a rating of: %d", game.White.Username, game.White.Rating)
	} else {
		log.Infof("Winner: %s with a rating of: %d", game.Black.Username, game.Black.Rating)
	}
	fmt.Println()
}

// getAllGamesForDate returns an empty array or an array of games for a given date.
func getAllGamesForDate(games []Game, day int, month time.Month, year int) []Game {
	var gamesWithGivenDate []Game
	for _, game := range games {
		gameEndTime := getTimeForUnixTimestamp(game.EndTime)
		if gameEndTime.Day() == day && gameEndTime.Month() == month && gameEndTime.Year() == year {
			gamesWithGivenDate = append(gamesWithGivenDate, game)
		}
	}
	return gamesWithGivenDate
}

// getTimeForUnixTimestamp is a simple wrapper for time.Unix() to safe some code while converting Unix Timestamps to Time Objects
func getTimeForUnixTimestamp(unixTimeStamp int) time.Time {
	return time.Unix(int64(unixTimeStamp), 0)
}

// Games contains an array of Game
type Games struct {
	Games []Game `json:"games"`
}

// Game holds all given Information about a Game no chessCom.com
type Game struct {
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
		Rating   int    `json:"rating"` // Represents the Rating after the game is finished
		Result   string `json:"result"` // Can be: "win", "resigned", "checkmated", "repetition", "stalemate", "timeout"
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
}
