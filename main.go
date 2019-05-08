package main

import (
	"encoding/json"
	"fmt"
	"github.com/nlopes/slack"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/receive", slashCommandHandler)
	fmt.Println("[INFO] Server listening")
	http.ListenAndServe(":"+os.Getenv("PORT"), nil)
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
		attachments := []interface{}{}
		attachments = append(attachments, map[string]interface{}{
			"text": " You asked me for " + params.Text,
		})
		response := map[string]interface{}{
			"response_type": "in_channel",
			"text":          "Hello <@" + s.UserID + ">",
			"attachments":   attachments,
		}
		data, err := json.Marshal(response)
		if err != nil {
			fmt.Errorf(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)

	default:
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
