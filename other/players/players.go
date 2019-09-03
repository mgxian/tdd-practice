package players

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type PlayerStore interface {
	GetPlayerScore(name string) int
	RecordWin(name string)
	GetLeaguePlayers() League
}

type PlayerServer struct {
	store PlayerStore
	http.Handler
}

func newPlayerServer(store PlayerStore) *PlayerServer {
	p := new(PlayerServer)
	p.store = store

	router := http.NewServeMux()
	router.HandleFunc("/players/", p.handlePlayer)
	router.HandleFunc("/league", p.handleLeague)

	p.Handler = router

	return p
}

func (p *PlayerServer) handlePlayer(w http.ResponseWriter, r *http.Request) {
	player := r.URL.Path[len("/players/"):]
	switch r.Method {
	case http.MethodPost:
		p.processWin(w, player)
	case http.MethodGet:
		p.showScore(w, player)
	}
}

const jsonContentType = "application/json"

func (p *PlayerServer) handleLeague(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", jsonContentType)
	leaguePlayers := p.getLeaguePlayers()
	json.NewEncoder(w).Encode(leaguePlayers)
}

func (p *PlayerServer) getLeaguePlayers() League {
	return p.store.GetLeaguePlayers()
}

func (p *PlayerServer) showScore(w http.ResponseWriter, player string) {
	score := p.store.GetPlayerScore(player)

	if score == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	fmt.Fprint(w, score)
}

func (p *PlayerServer) processWin(w http.ResponseWriter, player string) {
	w.WriteHeader(http.StatusAccepted)
	p.store.RecordWin(player)
}

type stubPlayerStore struct {
	scores   map[string]int
	winCalls []string
	league   League
}

func (s *stubPlayerStore) RecordWin(name string) {
	s.winCalls = append(s.winCalls, name)
}

func (s *stubPlayerStore) GetPlayerScore(name string) int {
	return s.scores[name]
}

func (s *stubPlayerStore) GetLeaguePlayers() League {
	return s.league
}

type inMemoryPlayerStore struct {
	store map[string]int
}

func newInMemoryPlayerStore() *inMemoryPlayerStore {
	return &inMemoryPlayerStore{
		store: make(map[string]int, 0),
	}
}

func (i *inMemoryPlayerStore) GetPlayerScore(name string) int {
	return i.store[name]
}

func (i *inMemoryPlayerStore) RecordWin(name string) {
	i.store[name]++
}

func (i *inMemoryPlayerStore) GetLeaguePlayers() League {
	var league League
	for p, w := range i.store {
		league = append(league, player{p, w})
	}
	return league
}

type player struct {
	Name string
	Wins int
}

type FileSystemStore struct {
	database io.ReadWriteSeeker
}

func (f *FileSystemStore) GetLeaguePlayers() League {
	f.database.Seek(0, 0)
	league, _ := NewLeague(f.database)
	return league
}

func (f *FileSystemStore) GetPlayerScore(name string) int {
	p := f.GetLeaguePlayers().Find(name)
	if p != nil {
		return p.Wins
	}
	return 0
}

func (f *FileSystemStore) RecordWin(name string) {
	league := f.GetLeaguePlayers()
	p := league.Find(name)
	if p != nil {
		p.Wins++
	} else {
		league = append(league, player{name, 1})
	}

	f.database.Seek(0, 0)
	json.NewEncoder(f.database).Encode(league)
}

type League []player

func (l League) Find(name string) *player {
	for i, p := range l {
		if p.Name == name {
			return &l[i]
		}
	}
	return nil
}

func NewLeague(r io.Reader) (League, error) {
	var league League
	err := json.NewDecoder(r).Decode(&league)
	if err != nil {
		err = fmt.Errorf("problem parsing league, %v", err)
	}
	return league, err
}
