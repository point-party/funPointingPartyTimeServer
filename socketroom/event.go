package socketroom

import (
	"encoding/json"
	"fmt"
	"log"
)

// determineGameAction takes in a gameMessage with a json rawMessage in.
// Using the event type we can determine how to decode it, and what action to take.
func determineGameAction(c *Client, gm *GameMessage, content json.RawMessage) {
	switch gm.Event {
	case voted:
		var ps PlayerStatus
		if err := json.Unmarshal(content, &ps); err != nil {
			log.Fatal(err)
		}
		c.Room.updateVote(ps.Name, ps.Point)
	default:
		fmt.Println("Did not match any known game events. No action taken.")
	}
}
