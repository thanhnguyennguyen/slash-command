package main

import (
	"encoding/json"
	"fmt"
	"github.com/nlopes/slack"
	"net/http"
	"regexp"
	"os"
	"strings"
)

// User struct which contains a name
type User struct {
	id string `json:"id"`
    name   string `json:"name"`
    manager string `json:"manager"`
}

var users map[string]*User

func main() {

	/**
	 Fetch user info
	**/

	users = map[string]*User {
		"U03RWP4LJDS" : {
			"U03RWP4LJDS",
			"nguyennguyen",
			"U03RWP4LJDS",
		},
	}
	/**
	END: Fetch user info
	**/


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
	regex := regexp.MustCompile("(0?[1-9]|[12][0-9]|3[01])/(0?[1-9]|1[012])/((19|20)\\d\\d)")

	switch s.Command {
	case "/off":
		params := &slack.Msg{Text: s.Text}

		attachments := []interface{}{}

		date := strings.Trim(params.Text, " ")
		if regex.MatchString(date) == false {
			attachments = append(attachments, map[string]interface{}{
				"text": "Invalid date format. DD/MM/YYYY",
			})
		} else {
			attachments = append(attachments, map[string]interface{}{
				"text": "<@" + s.UserID + "> is OFF on " + date + " . Message responses might be delayed.\n <@" + users[s.UserID].manager + "> should prepare options for backup if needed",
			})
		}

		response := map[string]interface{}{
			"response_type": "in_channel",
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

	case "/sick":
		params := &slack.Msg{Text: s.Text}

		attachments := []interface{}{}

		date := strings.Trim(params.Text, " ")
		if regex.MatchString(date) == false {
			attachments = append(attachments, map[string]interface{}{
				"text": "Invalid date format. DD/MM/YYYY",
			})
		} else {
			attachments = append(attachments, map[string]interface{}{
				"text": "<@" + s.UserID + "> is on MEDICAL LEAVE on " + date + " . Message responses might be delayed.\n <@" + users[s.UserID].manager + "> should prepare options for backup if needed",
			})
		}

		response := map[string]interface{}{
			"response_type": "in_channel",
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


	case "/remote":
		params := &slack.Msg{Text: s.Text}

		attachments := []interface{}{}

		date := strings.Trim(params.Text, " ")
		if regex.MatchString(date) == false {
			attachments = append(attachments, map[string]interface{}{
				"text": "Invalid date format. DD/MM/YYYY",
			})
		} else {
			attachments = append(attachments, map[string]interface{}{
				"text": "<@" + s.UserID + "> is WORKING REMOTELY on " + date + " . All working processes and communications should flow as normal.\nNote: Remote working are not encouraged in office-days.",
			})
		}

		response := map[string]interface{}{
			"response_type": "in_channel",
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


	case "/rules": 
		attachments := []interface{}{}
		attachments = append(attachments, map[string]interface{}{
			"text": "Office Rules:\n- Office days: Tuesday, Thursday, Friday\n- Remote days: Monday, Wednesday\n\nSynxtax:\n/off [date]: for paid leave, Example: 31/12/2023\n/remote [date]: for remote working leave, Example: 31/12/2023\n/sick [date]: for medical leave, Example: 31/12/2023",
		})
		response := map[string]interface{}{
			"response_type": "in_channel",
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
