package authentication

import (
	"fmt"
	"net/http"

	"golang.org/x/oauth2"
)

// GetToken contains logic to get turn code into token and response
func GetToken(config oauth2.Config, state string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		stateResp := r.FormValue("state")
		if stateResp != state {
			http.Error(w, "State code doesn't match", 400)
			return
		}
		code := r.FormValue("code")
		token, err := config.Exchange(oauth2.NoContext, code)
		if err != nil {
			fmt.Println("err", err)
			http.Error(w, "Couldn't get token from code", 500)
			return
		}
		resp := fmt.Sprintf("State: %s \n, Token: %s", stateResp, token)
		w.Write([]byte(resp))
	}
}
