package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/hilmihi/chirpbird/adapter"
	"github.com/rs/cors"
)

var listChannel []Channel = []Channel{
	{
		ID:           1,
		Name:         "Global Chat",
		Participants: 0,
		Sockets:      []string{},
	},
	{
		ID:           2,
		Name:         "Funny",
		Participants: 0,
		Sockets:      []string{},
	},
}

func main() {
	go h.run()

	port := "8080"

	mux := http.NewServeMux()

	//open db connection
	adapter.Open()

	if adapter.IsOpen() {
		log.Println("db: connection is already opened")
	}

	mux.Handle("/room",
		http.StripPrefix("/room", http.FileServer(http.Dir("./client"))),
	)

	mux.HandleFunc("/getChannels", func(w http.ResponseWriter, r *http.Request) {
		// output success response
		buf := new(bytes.Buffer)
		encoder := json.NewEncoder(buf)
		encoder.Encode(listChannel)
		w.Header().Set("Content-Type", "application/json")
		w.Write(buf.Bytes())
	})

	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		params, ok := r.URL.Query()["roomid"]
		paramsUsername, okUsername := r.URL.Query()["username"]

		if !ok || len(params[0]) < 1 {
			log.Println("Url Param 'key' is missing")
			return
		}

		if !okUsername || len(paramsUsername[0]) < 1 {
			log.Println("Url Param 'username' is missing")
			return
		}

		roomId := params[0]
		username := params[0]

		serveWs(w, r, roomId, username)
	})

	handler := cors.New(cors.Options{
		AllowedOrigins:   []string{},
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodDelete},
		AllowCredentials: true,
	}).Handler(mux)

	log.Print("Server starting at localhost:" + port)

	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatal(err)
	}
}

type Channel struct {
	ID           int64    `json:"id"`
	Name         string   `json:"name"`
	Participants int64    `json:"participants"`
	Sockets      []string `json:"sockets"`
	Channel_id   int64    `json:"channel_id"`
}

type ClientComMessage struct {
	ID           int64    `json:"id"`
	Channel_id   int64    `json:"channel_id"`
	Flag         string   `json:"flag"`
	Sendername   string   `json:"senderName"`
	Participants int64    `json:"participants"`
	Text         string   `json:"text"`
	Message      *Message `json:"message"`
	Channel      *Channel `json:"channel"`
}

type Message struct {
	ID         int64  `json:"id"`
	Channel_id int64  `json:"channel_id"`
	Sendername string `json:"senderName"`
	Content    string `json:"content"`
}
