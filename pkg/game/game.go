package game

import (
	"fmt"
	"sync"

	"github.com/dtimm/roomalone-plugin/pkg/session"

	"github.com/google/uuid"
)

type Engine struct {
	Sessions map[string]*session.Session
	*sync.Mutex
}

func New() *Engine {
	return &Engine{
		Sessions: map[string]*session.Session{},
		Mutex:    &sync.Mutex{},
	}
}

func (e *Engine) NewSession(game string) string {
	e.Lock()
	defer e.Unlock()

	g := session.New(game)
	u := uuid.New()

	e.Sessions[u.String()] = g

	return u.String()
}

func (e *Engine) GetSession(uuid string) (*session.Session, error) {
	e.Lock()
	defer e.Unlock()

	if v, ok := e.Sessions[uuid]; ok {
		l := v.TryLock()
		if !l {
			return nil, fmt.Errorf("session %s currently in use", uuid)
		}
		v.Unlock()
		return v, nil
	}

	return nil, fmt.Errorf("could not find session ID %s", uuid)
}

func (e *Engine) EndSession(uuid string) error {
	e.Lock()
	defer e.Unlock()

	if _, ok := e.Sessions[uuid]; ok {
		delete(e.Sessions, uuid)
		return nil
	}

	return fmt.Errorf("could not find session ID %s", uuid)
}
