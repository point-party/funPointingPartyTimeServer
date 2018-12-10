package socketroom

import (
	"encoding/json"
	"fmt"
	"log"
)

// determineGameAction takes in a gameMessage with a json rawMessage in.
// Using the event type we can determine how to decode it, and what action to take.
func determineGameAction(r *Room, gm GameMessage, content json.RawMessage) {
	switch gm.Event {
	case voted:
		var ps PlayerStatus
		if err := json.Unmarshal(content, &ps); err != nil {
			log.Fatal(err)
		}
		r.updateVote(ps.Point, ps.ID)
	case clearPoints:
		r.clearPoints()
	case revealPoints:
		fmt.Println("In reveal case")
	default:
		fmt.Println("Did not match any known game events. No action taken.")
	}
}
