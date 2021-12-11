package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	guuid "github.com/google/uuid"
	"github.com/hilmihi/chirpbird/adapter"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

func main() {
	go h.run()

	var err error
	err = godotenv.Load()
	if err != nil {
		log.Fatalf("Error getting env, not comming through %v", err)
	} else {
		log.Println("We are getting the env values")
	}

	port := "8080"

	mux := http.NewServeMux()

	//open db connection
	err = adapter.Open(os.Getenv("DB_DRIVER"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_PORT"), os.Getenv("DB_HOST"), os.Getenv("DB_NAME"))

	if err != nil {
		log.Print(err)
	} else {
		log.Println("db: connection is already opened")
	}

	mux.HandleFunc("/getChannels", func(w http.ResponseWriter, r *http.Request) {
		//get username from param URL
		paramsUsername, okUsername := r.URL.Query()["username"]

		if !okUsername || len(paramsUsername[0]) < 1 {
			log.Println("Url Param 'username' is missing")
			return
		}

		//get data user from db based on param username
		user, err := adapter.UserGet(paramsUsername[0])
		if err != nil {
			log.Println(err)
		}

		//get all room that user already subscribed
		subs, err := adapter.SubsByUser(user.ID)
		if err != nil {
			log.Println(err)
		}

		// output success response
		buf := new(bytes.Buffer)
		encoder := json.NewEncoder(buf)
		encoder.Encode(subs)
		w.Header().Set("Content-Type", "application/json")
		w.Write(buf.Bytes())
	})

	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		// Read body
		b, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		// Unmarshal
		var msg ClientComMessage
		err = json.Unmarshal(b, &msg)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		user, err := adapter.UserGet(msg.Username)
		if err != nil {
			log.Println(err)
		}

		if user == nil {
			idCreate, err := adapter.UserCreate(&adapter.User{
				CreateDate: time.Now().UTC().Round(time.Millisecond),
				UpdateDate: time.Now().UTC().Round(time.Millisecond),
				Username:   msg.Username,
				State:      1,
				Lastseen:   time.Now().UTC().Round(time.Millisecond),
			})

			if err != nil {
				log.Println(err)
				return
			}
			user, err = adapter.UserGetByID(idCreate)
			if err != nil {
				log.Println(err)
				return
			}
		}

		// output success response
		buf := new(bytes.Buffer)
		encoder := json.NewEncoder(buf)
		encoder.Encode(user)
		w.Header().Set("Content-Type", "application/json")
		w.Write(buf.Bytes())
	})

	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		//get username from param URL
		paramsUsername, okUsername := r.URL.Query()["username"]

		if !okUsername || len(paramsUsername[0]) < 1 {
			log.Println("Url Param 'username' is missing")
			return
		}

		//get data user from db based on param username
		user, err := adapter.UserGet(paramsUsername[0])
		if err != nil {
			log.Println(err)
		}

		//get all room that user already subscribed
		subs, err := adapter.SubsByUser(user.ID)
		if err != nil {
			log.Println(err)
		}

		roomId := guuid.New().String()

		serveWs(w, r, roomId, subs)
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

type ClientComMessage struct {
	ID            int64                  `json:"id"`
	Channel_id    int64                  `json:"channel_id"`
	Flag          string                 `json:"flag"`
	Username      string                 `json:"username"`
	Text          string                 `json:"text"`
	MessageB      []adapter.Message      `json:"messageb"`
	MessageC      *Message               `json:"messagec"`
	Room          *adapter.Room          `json:"room"`
	Users         []adapter.User         `json:"users"`
	Subscriptions []adapter.Subscription `json:"subcriptions"`
}

type Message struct {
	ID         int64  `json:"id"`
	Channel_id int64  `json:"channel_id"`
	Username   string `json:"username"`
	Content    string `json:"Content"`
}
