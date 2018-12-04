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
		vote := voteEvent{ps.Point, ps.ID}
		c.Room.vote <- vote
	case clearPoints:
		c.Room.clear <- true
	case revealPoints:
		fmt.Println("In reveal case")
	default:
		fmt.Println("Did not match any known game events. No action taken.")
	}
}
