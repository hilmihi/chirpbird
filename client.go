package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	guuid "github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/hilmihi/chirpbird/adapter"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var c *connection

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(*http.Request) bool { return true },
}

// connection is an middleman between the websocket connection and the hub.
type connection struct {
	// The websocket connection.
	ws *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte
}

// readPump pumps messages from the websocket connection to the hub.
func (s subscription) readPump() {
	c := s.conn
	defer func() {
		h.unregister <- s
		c.ws.Close()
	}()
	c.ws.SetReadLimit(maxMessageSize)
	c.ws.SetReadDeadline(time.Now().Add(pongWait))
	c.ws.SetPongHandler(func(string) error { c.ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, msg, err := c.ws.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("error: %v", err)
			}
			break
		}
		s.dispatchRaw(msg)
	}
}

// write writes a message with the given message type and payload.
func (c *connection) write(mt int, payload []byte) error {
	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return c.ws.WriteMessage(mt, payload)
}

// writePump pumps messages from the hub to the websocket connection.
func (s *subscription) writePump() {
	c := s.conn
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.ws.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.write(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.write(websocket.TextMessage, message); err != nil {
				return
			}
		case <-ticker.C:
			if err := c.write(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

func (s *subscription) dispatchRaw(raw []byte) {
	var msg ClientComMessage

	if err := json.Unmarshal(raw, &msg); err != nil {
		// Malformed message
		log.Fatal(err)
		return
	}

	s.dispatch(&msg)
}

func (s *subscription) dispatch(msg *ClientComMessage) {
	switch {
	case msg.Flag == "channel-join":
		msgs, err := adapter.MessageByRoom(msg.Room.Name)
		if err != nil {
			log.Println(err)
		}

		msg.MessageB = msgs

		msgBc, err := json.Marshal(msg)
		if err != nil {
			log.Println(err.Error())
			return
		}
		m := message{msgBc, s.room}
		h.broadcast <- m
	case msg.Flag == "message":
		var status_save bool = false
		msgs, err := adapter.MessageByRoom(msg.Room.Name)
		if err != nil {
			log.Println(err)
		} else {
			var seq int = len(msgs) + 1
			messag := &adapter.Message{
				CreateDate: time.Now().UTC().Round(time.Millisecond),
				UpdateDate: time.Now().UTC().Round(time.Millisecond),
				Seqid:      int64(seq),
				Room:       msg.Room.Name,
				From:       msg.ID,
				Content:    msg.MessageC.Content,
			}
			err := adapter.SaveMessage(messag)

			if err == nil {
				status_save = true
			}
		}

		if status_save {
			msgBc, err := json.Marshal(msg)
			if err != nil {
				log.Println(err.Error())
				return
			}
			m := message{msgBc, msg.Room.Name}
			h.broadcast <- m
		}
	case msg.Flag == "find-user":
		users, err := adapter.FindUsers(msg.Text, msg.ID)
		if err != nil {
			log.Println(err)
		}

		msg.Users = users

		msgBc, err := json.Marshal(msg)
		if err != nil {
			log.Println(err.Error())
			return
		}
		m := message{msgBc, s.room}
		h.broadcast <- m
	case msg.Flag == "create-room":
		initiator := &adapter.Subscription{
			CreateDate: time.Now().UTC().Round(time.Millisecond),
			UpdateDate: time.Now().UTC().Round(time.Millisecond),
			UserID:     msg.ID,
		}

		room := &adapter.Room{
			CreateDate: time.Now().UTC().Round(time.Millisecond),
			UpdateDate: time.Now().UTC().Round(time.Millisecond),
			Name:       guuid.New().String(),
			Public:     msg.Room.Public,
		}

		rm, err := adapter.RoomCreate(room)
		if err != nil {
			log.Println(err)
			return
		}
		room, err = adapter.RoomGetByID(rm)

		initiator.Room = room.Name
		err = adapter.RoomCreateP2P(initiator, msg.Users)
		if err != nil {
			log.Println(err)
			return
		}

		subs, err := adapter.SubsByUser(msg.ID)
		if err != nil {
			log.Println(err)
		}

		msg.Subscriptions = subs
		msg.Flag = "get-channel"
		msg.Room = nil
		msg.Users = nil

		msgBc, err := json.Marshal(msg)
		if err != nil {
			log.Println(err.Error())
			return
		}
		m := message{msgBc, s.room}
		h.broadcast <- m
	default:
		// Unknown message
		log.Print("Unknown message")
		return
	}
}

// serveWs handles websocket requests from the peer.
func serveWs(w http.ResponseWriter, r *http.Request, roomId string, subs []adapter.Subscription) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err.Error())
		return
	}
	c = &connection{send: make(chan []byte, 256), ws: ws}

	s := subscription{c, roomId}
	h.register <- s
	go s.writePump()
	go s.readPump()

	subsChannel(subs)
}

func subsChannel(subs []adapter.Subscription) {
	for _, sub := range subs {
		s := subscription{c, sub.Room}
		h.register <- s
	}
}
