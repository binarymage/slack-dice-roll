package main

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"regexp"
	"strconv"
	"time"
	"fmt"
)

type SlackMessage struct {
	Text string
	Response_type string
}

var dice *regexp.Regexp = regexp.MustCompile(`^(\d*)d(\d+)$`)

func send_json(w http.ResponseWriter, text string, in_channel bool) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)

	message := SlackMessage {
		Text: text,
	}

	if in_channel {
		message.Response_type = "in_channel"
	} else {
		message.Response_type = "ephemeral"
	}

	if b, err := json.Marshal(message); err == nil {
		w.Write(b)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		send_json(w, "Bad HTTP Method", false)
		return
	}

	if token := r.FormValue("token"); token != "change_me" {
		send_json(w, "Bad token", false)
		return
	}

	if cmd := r.FormValue("command"); cmd != "/roll" {
		send_json(w, "Invalid command " + cmd, false)
		return
	}

	text := r.FormValue("text")

	var num_dice, num_faces int = 1, 6

	if dice.MatchString(text) {
		matches := dice.FindStringSubmatch(text)
		if matches[1] != "" {
			num_dice, _ = strconv.Atoi(matches[1])
		}
		num_faces, _ = strconv.Atoi(matches[2])
	}

	user_id := r.FormValue("user_id")
	user_name := r.FormValue("user_name")
	rand.Seed(time.Now().UnixNano())

	// roll dice
	var sum int = 0

	for i := 0; i < num_dice; i++ {
		sum += (rand.Intn(num_faces) + 1)
	}

	// "<@" + user_id + "|"+ user_name +"> rolled " + sum + " from " + num_dice + " d" + num_faces
	message := fmt.Sprintf("<@%s|%s> rolled %d from %dd%d", user_id, user_name, sum, num_dice, num_faces)
	send_json(w, message, false)
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":9000", nil)
}