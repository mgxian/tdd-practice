package players

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type PlayerStore interface {
	GetPlayerScore(name string) int
	RecordWin(name string)
	GetLeaguePlayers() []player
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

func (p *PlayerServer) getLeaguePlayers() []player {
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
	league   []player
}

func (s *stubPlayerStore) RecordWin(name string) {
	s.winCalls = append(s.winCalls, name)
}

func (s *stubPlayerStore) GetPlayerScore(name string) int {
	return s.scores[name]
}

func (s *stubPlayerStore) GetLeaguePlayers() []player {
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

func (s *inMemoryPlayerStore) GetPlayerScore(name string) int {
	return s.store[name]
}

func (s *inMemoryPlayerStore) RecordWin(name string) {
	s.store[name]++
}

func (s *inMemoryPlayerStore) GetLeaguePlayers() []player {
	return nil
}

type player struct {
	Name string
	Wins int
}
