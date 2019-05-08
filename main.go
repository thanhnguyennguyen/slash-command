package main

import (
	"fmt"
	"github.com/nlopes/slack"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/receive", slashCommandHandler)
	fmt.Println("[INFO] Server listening")
	http.ListenAndServe(":" + os.Getenv("PORT"), nil)
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
		response := fmt.Sprintf(s.UserID + ", you asked for the me for %v", params.Text)
		w.Write([]byte(response))

	default:
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
