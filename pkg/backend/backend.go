package backend

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/dtimm/roomalone-plugin/pkg/game"
	"github.com/dtimm/roomalone-plugin/pkg/session"
	"github.com/gorilla/mux"
)

type GameBackend struct {
	*game.Engine
	PrintDebug bool
}

type handlerFuncWithSession func(w http.ResponseWriter, r *http.Request, s *session.Session)

func (g *GameBackend) HandleWithSession(h handlerFuncWithSession) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		guid := vars["session_guid"]

		if g.PrintDebug {
			fmt.Printf("req url: %s\n", r.URL)
			fmt.Printf("req method: %s\n", r.Method)
			fmt.Printf("req body: %s\n", r.Body)
		}

		s, err := g.GetSession(guid)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error getting session: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		h(w, r, s)
	}
}

func (g *GameBackend) Debug(w http.ResponseWriter, r *http.Request, s *session.Session) {
	b, err := json.Marshal(s)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error marshalling response body: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(b)
}

func (g *GameBackend) Inventory(w http.ResponseWriter, r *http.Request, s *session.Session) {
	if r.Method == "POST" {
		defer r.Body.Close()
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error reading request body: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		i := session.Inventory{}
		err = json.Unmarshal(body, &i)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error unmarshalling request body: %s, body: %s", err, body)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		s.SetInventory(i.Items)
	}

	b, err := json.Marshal(s.GetInventory())
	if err != nil {
		fmt.Fprintf(os.Stderr, "error marshalling response body: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if g.PrintDebug {
		fmt.Printf("inventory: %s\n", string(b))
	}

	w.Write(b)
}

type MoveRequest struct {
	Location string `json:"location"`
}

func (g *GameBackend) Location(w http.ResponseWriter, r *http.Request, s *session.Session) {
	var err error
	var l session.Location

	switch r.Method {
	case http.MethodGet:
		l, err = s.GetLocation("")
		if err != nil {
			fmt.Fprintf(os.Stderr, "error getting location: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	case http.MethodPatch:
		defer r.Body.Close()
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error reading request body: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		req := session.Location{}
		err = json.Unmarshal(body, &req)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error unmarshalling request body: %s, body: %s", err, body)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		l = s.SetLocation(req)
	case http.MethodPost:
		defer r.Body.Close()
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error reading request body: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		req := MoveRequest{}
		err = json.Unmarshal(body, &req)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error unmarshalling request body: %s, body: %s", err, body)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		l, err = s.MoveLocation(req.Location)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error setting location: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	b, err := json.Marshal(l)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error marshalling response body: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if g.PrintDebug {
		fmt.Printf("location: %s\n", string(b))
	}

	w.Write(b)
}

func (g *GameBackend) HandleEnd(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	guid := vars["session_guid"]

	err := g.EndSession(guid)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error getting session: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if g.PrintDebug {
		fmt.Printf("session ended: %s\n", guid)
	}

	w.WriteHeader(http.StatusOK)
}

type NewSessionRequest struct {
	Adventure string `json:"adventure"`
}

type NewSessionResponse struct {
	SessionGUID string `json:"session_guid"`
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

	if g.PrintDebug {
		fmt.Printf("new session: %s\n", string(b))
	}

	w.Write(b)
}

func New(g *game.Engine) *GameBackend {
	return &GameBackend{
		Engine: g,
	}
}
