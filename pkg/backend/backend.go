package backend

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/dtimm/roomalone-plugin/pkg/game"
	"github.com/gorilla/mux"
)

type GameBackend struct {
	*game.Engine
}

type NewSessionRequest struct {
	Adventure string `json:"adventure"`
}

type NewSessionResponse struct {
	SessionGUID string `json:"session_guid"`
}

func (g *GameBackend) HandleLocation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	guid := vars["session_guid"]

	s, err := g.GetSession(guid)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error getting session: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	l, err := s.GetLocation("")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error getting location: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(l)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error marshalling response body: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(b)
}

func (g *GameBackend) HandleNewSession(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading request body: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var n NewSessionRequest
	err = json.Unmarshal(body, &n)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error unmarshalling request body: %s, body: %s", err, body)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	s := g.NewSession(fmt.Sprintf("story/%s", n.Adventure))

	b, err := json.Marshal(NewSessionResponse{
		SessionGUID: s,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "error marshalling response body: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(b)
}

func New(g *game.Engine) *GameBackend {
	return &GameBackend{
		Engine: g,
	}
}
