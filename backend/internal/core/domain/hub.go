package domain

import "fmt"

type Room struct {
	Id      int             `json:"id"`
	Name    string          `json:"name"`
	Clients map[int]*Client `json:"clients"`
}

type PrivateRoom struct {
	Id      int             `json:"id"`
	Name    string          `json:"name"`
	Clients map[int]*Client `json:"clients"`
}

type Hub struct {
	Rooms        map[int]*Room
	PrivateRooms map[int]*PrivateRoom
	Register     chan *Client
	Unregister   chan *Client
	Broadcast    chan *Chat
	// Notify chan *Message
}

func NewHub() *Hub { // move it to service later
	return &Hub{
		Rooms:        make(map[int]*Room),
		PrivateRooms: make(map[int]*PrivateRoom),
		Register:     make(chan *Client),
		Unregister:   make(chan *Client),
		Broadcast:    make(chan *Chat, 1024),
		// Notify:  make(chan *Message, 1024),

	}
}

func (h *Hub) Run() {
	for {
		select {
		case cl := <-h.Register:

			//check the type of the client: group or private and add to the corresponding room
			if cl.Type == "group" {
				if _, ok := h.Rooms[cl.RoomId]; ok {
					r := h.Rooms[cl.RoomId]
					if _, ok := r.Clients[cl.Id]; !ok {
						fmt.Println("Register => ", cl)
						r.Clients[cl.Id] = cl
					}
				}
			} else {
				if _, ok := h.PrivateRooms[cl.RoomId]; ok {
					r := h.PrivateRooms[cl.RoomId]
					if _, ok := r.Clients[cl.Id]; !ok {
						r.Clients[cl.Id] = cl
					}
				}
			}
		case cl := <-h.Unregister:
			//check the type of the client: group or private and remove from the corresponding room
			// fmt.Println("Unregister => ", cl)
			if cl.Type == "group" {
				if _, ok := h.Rooms[cl.RoomId]; ok {
					if _, ok := h.Rooms[cl.RoomId].Clients[cl.Id]; ok {
						if len(h.Rooms[cl.RoomId].Clients) != 0 {
							h.Broadcast <- &Chat{
								Content:  "user left the chat",
								GroupId:  cl.RoomId,
								Username: cl.Username,
							}
						}
						delete(h.Rooms[cl.RoomId].Clients, cl.Id)
						close(cl.Message)
					}
				}
			} else {
				if _, ok := h.PrivateRooms[cl.RoomId]; ok {
					if _, ok := h.PrivateRooms[cl.RoomId].Clients[cl.Id]; ok {
						if len(h.PrivateRooms[cl.RoomId].Clients) != 0 {
							h.Broadcast <- &Chat{
								Content:  "user left the chat",
								GroupId:  cl.RoomId,
								Username: cl.Username,
							}
						}
						delete(h.PrivateRooms[cl.RoomId].Clients, cl.Id)
						close(cl.Message)
					}
				}
			}
		case m := <-h.Broadcast:
			//check the type of the msg: group or private and send into the correct room
			if m.Type == "group" {
				if roo, ok := h.Rooms[m.GroupId]; ok {
					for _, cl := range h.Rooms[m.GroupId].Clients {
						fmt.Println("Broadcast", roo)
						cl.Message <- m
					}
				}
			} else {
				fmt.Println("Sending private message to client", m)
				if _, ok := h.PrivateRooms[m.GroupId]; ok {
					fmt.Println("Room found", h.PrivateRooms[m.GroupId])
					fmt.Println("Clients found", h.PrivateRooms[m.GroupId].Clients)
					for _, cl := range h.PrivateRooms[m.GroupId].Clients {
						cl.Message <- m
					}
				}
			}

		}
	}
}
