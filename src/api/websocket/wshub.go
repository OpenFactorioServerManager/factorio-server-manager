package websocket

import (
	"reflect"

	"github.com/OpenFactorioServerManager/factorio-server-manager/bootstrap"
)

// the hub, that is exported and can be used anywhere to work with the websocket
var WebsocketHub *wsHub

var LogCache []string

// a controlHandler is used to determine, if something has to be done, on a specific command.
// register a handler with `wsHub.RegisterControlHandler`
// unregister a handler with `wsHub.UnregisterControlHandler`
type controlHandler func(controls WsControls)

// The type for of control messages.
// Type ans Value both have to be set, if controls are sent.
// Currently supported Type:
// - `subscribe` - Value used as room name
// - `unsubscribe` - Value used as room name
// - `command` - Value contains the command to execute
type WsControls struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

// the main message of our websocket protocol.
// if the room_name is an empty string, this will be sent as broadcast
// if controls is not empty, this will not be sent anywhere, but used as commands, to join/leave rooms and to send commands to the factorio server
type wsMessage struct {
	RoomName string      `json:"room_name"`
	Message  interface{} `json:"message,omitempty"`
	Controls WsControls  `json:"controls,omitempty"`
}

type wsRoom struct {
	// same as the key of the map in the wsHub
	name string

	// clients that are in this room. This list is a sublist of the one inside the hub
	clients map[*wsClient]bool

	// register a client to this room
	register chan *wsClient

	// unregister a client from this room, if no clients remain, this room will be deleted
	unregister chan *wsClient

	// send a message to all clients in this room
	send chan wsMessage
}

// Hub is the basic setup of the server.
// It contains everything needed for the websocket to run.
// Only the controlHandler Subscriptions are public, everything else can be controlled with the functions and the wsClient.
type wsHub struct {
	// list of all connected clients
	clients map[*wsClient]bool

	// Messages that should be sent to ALL clients
	broadcast chan wsMessage

	// a list of all rooms
	rooms map[string]*wsRoom

	// register a client to this hub
	register chan *wsClient

	// unregister a client from this hub
	unregister chan *wsClient

	// run a control message on all registered controlHandler
	runControl chan WsControls

	// list of all registered controlHandlers
	controlHandlers map[reflect.Value]controlHandler

	// register a controlHandler
	RegisterControlHandler chan controlHandler

	// unregister a controlHandler
	UnregisterControlHandler chan controlHandler
}

// initialize and run the mein websocket hub.
func init() {
	WebsocketHub = &wsHub{
		broadcast:                make(chan wsMessage),
		register:                 make(chan *wsClient),
		rooms:                    make(map[string]*wsRoom),
		unregister:               make(chan *wsClient),
		clients:                  make(map[*wsClient]bool),
		runControl:               make(chan WsControls),
		controlHandlers:          make(map[reflect.Value]controlHandler),
		RegisterControlHandler:   make(chan controlHandler),
		UnregisterControlHandler: make(chan controlHandler),
	}

	go WebsocketHub.run()
}

// remove a client from this hub and all of its rooms
func (hub *wsHub) removeClient(client *wsClient) {
	delete(hub.clients, client)
	close(client.send)
	for _, room := range hub.rooms {
		room.unregister <- client
	}
}

// run starts a websocket hub, this has to be done in a subroutine `go hub.run()`
func (hub *wsHub) run() {
	for {
		select {
		case client := <-hub.register:
			hub.clients[client] = true
		case client := <-hub.unregister:
			if _, ok := hub.clients[client]; ok {
				hub.removeClient(client)
			}
		case message := <-hub.broadcast:
			for client := range hub.clients {
				select {
				case client.send <- message:
				default:
					hub.removeClient(client)
				}
			}
		case function := <-hub.RegisterControlHandler:
			hub.controlHandlers[reflect.ValueOf(function)] = function
		case function := <-hub.UnregisterControlHandler:
			delete(hub.controlHandlers, reflect.ValueOf(function))
		}
	}
}

// Broadcast a message to all connected clients (only clients connected to this room).
func (hub *wsHub) Broadcast(message interface{}) {
	hub.broadcast <- wsMessage{
		RoomName: "",
		Message:  message,
	}
}

// get a websocket room or create it, if it doesn't exist yet.
// Also starts the rooms subroutine `wsRoom.run()`
func (hub *wsHub) GetRoom(name string) *wsRoom {
	if room, ok := hub.rooms[name]; ok {
		return room
	} else {
		room := &wsRoom{
			name:       name,
			clients:    make(map[*wsClient]bool),
			register:   make(chan *wsClient),
			unregister: make(chan *wsClient),
			send:       make(chan wsMessage),
		}
		hub.rooms[name] = room
		go room.run()
		return room
	}
}

// run starts a websocket room. This has to be run as a subroutine `go room.run()`
func (room *wsRoom) run() {
	for {
		select {
		case client := <-room.register:
			room.clients[client] = true

			// some hardcoded stuff for gamelog room
			if room.name == "gamelog" {
				// send cached log to registered client
				for _, logLine := range LogCache {
					client.send <- wsMessage{
						RoomName: "gamelog",
						Message:  logLine,
					}
				}
			}
		case client := <-room.unregister:
			if _, ok := room.clients[client]; ok {
				delete(room.clients, client)
				// FIXME when more rooms are used, remove empty rooms.
				// Since we only have a few rooms at the same time, just keep them.
				// This is code, that will cause a concurrent call on `wsHub.rooms`.
				// To fix this, move the deletion into the hub.
				// Be careful to think about race conditions, if a user registered to the room, before room was really deleted.
				//if len(room.clients) == 0 {
				//	//remove this room
				//	delete(room.hub.rooms, room.name)
				//	return
				//}
			}
		case message := <-room.send:
			for client := range room.clients {
				select {
				case client.send <- message:
				default:
					room.unregister <- client
				}
			}

			// some hardcoded stuff for gamelog room
			if room.name == "gamelog" {
				// add the line to the cache
				LogCache = append(LogCache, message.Message.(string))
				config := bootstrap.GetConfig()

				// When cache is bigger than max size, delete one line
				if len(LogCache) > config.ConsoleCacheSize {
					LogCache = LogCache[1:]
				}
			}
		}
	}
}

// Send a message into this room.
func (room *wsRoom) Send(message interface{}) {
	room.send <- wsMessage{
		RoomName: room.name,
		Message:  message,
	}
}
