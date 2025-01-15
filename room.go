package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

// ClientMessage contains the message and which client sent it
type ClientMessage struct {
	client *Client
	data   []byte
}

// GameRoom maintains the set of active clients and broadcasts messages to the
// clients.
type GameRoom struct {
	// Registered clients, holds the ws connection to the clients.
	clients map[*Client]bool

	// Channel for broadcasting game state (positions, velocities etc)
	output chan []byte

	// Channel for receiving player input
	input chan ClientMessage

	// Register requests from the clients
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	// each player will be mapped to the corresponding connection
	players map[*Client]*Player
}

func NewRoom() *GameRoom {
	return &GameRoom{
		input:      make(chan ClientMessage),
		output:     make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		players:    make(map[*Client]*Player),
	}
}

func (room *GameRoom) run() {
	for {

		for _, player := range room.players {
			player.step()
		}

		select {
		case client := <-room.register:
			newPlayer := &Player{
				ID:           client.id,
				Position:     Vec2{0, 0},
				Velocity:     Vec2{0, 0},
				Acceleration: Vec2{0, 0},

				input: NewInputState(),
			}

			room.clients[client] = true
			room.players[client] = newPlayer
		case client := <-room.unregister:
			if _, ok := room.clients[client]; ok {
				delete(room.clients, client)
				close(client.send)
			}

		case message := <-room.input:
			client := message.client
			data := string(message.data)
			dataSplit := strings.Split(data, "-")

			fmt.Println(data)
			if strings.Contains(dataSplit[0], "pressed") {
				if strings.Contains(dataSplit[1], "W") {
					room.players[client].input.Update(UP, PRESSED)
				}

				if dataSplit[1] == "S" {
					room.players[client].input.Update(DOWN, PRESSED)
				}

				if data[1] == 'A' {
					room.players[client].input.Update(DOWN, PRESSED)
				}

				if data[1] == 'D' {
					room.players[client].input.Update(RIGHT, PRESSED)
				}
			} else {
				if data[1] == 'W' {
					room.players[client].input.Update(UP, UNPRESSED)
				}

				if data[1] == 'S' {
					room.players[client].input.Update(DOWN, UNPRESSED)
				}

				if data[1] == 'A' {
					room.players[client].input.Update(DOWN, UNPRESSED)
				}

				if data[1] == 'D' {
					room.players[client].input.Update(RIGHT, UNPRESSED)
				}
			}

		case message := <-room.output:
			for client := range room.clients {
				client.send <- message
			}
		}
	}
}

func (room *GameRoom) broadcast() {
	var state []byte

	for _, player := range room.players {
		byteSlice, err := json.Marshal(player)
		if err != nil {
			fmt.Println("JSON marshal failed:", err)
			return
		}

		state = append(state, byteSlice...)
	}

	room.output <- state
}
