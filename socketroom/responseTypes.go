package socketroom

// RoomName contains the name in json format to send to the client.
type RoomName struct {
	Name string `json:"roomName"`
}

// GameMessage will be the json structure used to communicate
type GameMessage struct {
	Event   string      `json:"event"`
	Payload interface{} `json:"payload"`
}

// PlayerUpdate contains a list of players and their current point.
type PlayerUpdate struct {
	Players []PlayerStatus `json:"players"`
}

// PlayerStatus has list of players and their current points.
type PlayerStatus struct {
	Name  string `json:"name"`
	Point string `json:"point"`
	ID    string `json:"id"`
}

// GAME CONSTANTS
const (
	joinRoom     = "JOIN_ROOM"
	leaveRoom    = "LEAVE_ROOM"
	submitPoint  = "SUBMIT_POINT"
	revealPoints = "REVEAL_POINTS"
	clearPoints  = "CLEAR_POINTS"
	voted        = "VOTED"
)
