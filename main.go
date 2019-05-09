package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/nlopes/slack"
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
		break

	case "/luckyroll":
		if strings.ToLower(s.Text) == "lunch" {
			apiKey := os.Getenv("API_KEY")
			// search message to get thread_ts
			url := "https://slack.com/api/search.messages?token" + apiKey + "&query='Lunch tomorrow at the CIC building'"

			req, err := http.NewRequest("GET", url, bytes.NewBuffer([]byte{}))

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			defer resp.Body.Close()

			body, _ := ioutil.ReadAll(resp.Body)
			var searchMessages slack.SearchMessages
			json.Unmarshal(body, &searchMessages)
			channel := searchMessages.Matches[0].Channel
			ts := searchMessages.Matches[0].Timestamp
			url = "https: //slack.com/api/channels.replies?token" + apiKey + "&channel=" + channel + "&thread_ts=" + ts
			req, err = http.NewRequest("GET", url, bytes.NewBuffer([]byte{}))

			client = &http.Client{}
			resp, err = client.Do(req)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			defer resp.Body.Close()

			body, _ = ioutil.ReadAll(resp.Body)
			var ms []slack.GetConversationHistoryResponse
			json.Unmarshal(body, &ms)
			repliers := ms[0].ReplyUsers
			// get replier list
			// random
		}

		break

	default:
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
