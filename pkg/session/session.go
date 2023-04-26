package session

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"sync"
)

type Session struct {
	CurrentLocation string
	Locations       map[string]Location
	Inventory
	*sync.Mutex
}

type Inventory struct {
	Items []string `json:"items"`
}

type Location struct {
	Name        string   `json:"name"`
	Connections []string `json:"connections"`
	Description string   `json:"description"`
	Changes     []string `json:"changes"`
	Story       string   `json:"story"`
}

func New(game string) *Session {
	l := loadLocations(fmt.Sprintf("%s/locations.json", game))
	i := loadInventory(fmt.Sprintf("%s/inventory.json", game))
	return &Session{
		Inventory:       i,
		Locations:       l,
		CurrentLocation: "start-location",
		Mutex:           &sync.Mutex{},
	}
}

func loadInventory(f string) Inventory {
	i, err := os.ReadFile(f)
	if err != nil {
		panic(err)
	}

	ir := new(Inventory)
	err = json.Unmarshal(i, ir)
	if err != nil {
		panic(err)
	}

	return *ir
}

func loadLocations(f string) map[string]Location {
	contents, err := os.ReadFile(f)
	if err != nil {
		panic(err)
	}

	l := new(map[string]Location)
	err = json.Unmarshal(contents, l)
	if err != nil {
		panic(err)
	}

	return *l
}

func (s *Session) SetLocation(input string) (Location, error) {
	s.Lock()
	defer s.Unlock()

	if l, ok := s.Locations[input]; ok {
		s.CurrentLocation = input
		return l, nil
	}

	return Location{}, fmt.Errorf("location %s not in map: Move action only accepts location names from the current location connections", input)
}

func (s *Session) GetInventory() Inventory {
	s.Lock()
	defer s.Unlock()

	return s.Inventory
}

func (s *Session) SetInventory(i []string) {
	s.Lock()
	defer s.Unlock()

	s.Inventory = Inventory{Items: i}
}

func (s *Session) GetLocation(input string) (Location, error) {
	s.Lock()
	defer s.Unlock()

	if input != "" {
		if _, ok := s.Locations[input]; !ok {
			return Location{}, fmt.Errorf("location %s not in map: Location action only accepts existing location names or empty input", input)
		}
	} else {
		input = s.CurrentLocation
	}

	return s.Locations[input], nil
}

func (s *Session) AddLocationChange(input string) (Location, error) {
	re := regexp.MustCompile(`^([a-z0-9\-]+): (.*)$`)
	m := re.FindAllStringSubmatch(input, 1)
	if len(m) != 1 || len(m[0]) != 3 {
		return Location{}, fmt.Errorf("input `%s` did not match format: `location: description of change`", input)
	}

	match := m[0]
	loc := match[1]

	s.Lock()
	if v, ok := s.Locations[loc]; ok {
		v.Changes = append(s.Locations[loc].Changes, match[2])
		s.Locations[loc] = v

		s.Unlock()
		return s.GetLocation(loc)
	}

	return Location{}, fmt.Errorf("location %s not in map: only use location names that you know exist", loc)
}
