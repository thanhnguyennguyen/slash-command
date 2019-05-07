package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/nlopes/slack"
	"log"
	"net/http"
	"os"
)

func main() {
	err := godotenv.Load("env.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	http.HandleFunc("/receive", slashCommandHandler)

	fmt.Println("[INFO] Server listening")
	http.ListenAndServe(":8080", nil)
}

func slashCommandHandler(w http.ResponseWriter, r *http.Request) {
	s, err := slack.SlashCommandParse(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !s.ValidateToken(os.Getenv("SLACK_VERIFICATION_TOKEN")) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	switch s.Command {
	case "/testbot":
		params := &slack.Msg{Text: s.Text}
		response := fmt.Sprintf("You asked for the me for %v", params.Text)
		w.Write([]byte(response))

	default:
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
