package players

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
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
	database *json.Encoder
	league   League
}

func initializePlayerDBFile(database *os.File) error {
	database.Seek(0, 0)

	info, err := database.Stat()
	if err != nil {
		return fmt.Errorf("problem getting file info from file %s, %v", info.Name(), err)
	}

	if info.Size() == 0 {
		database.WriteString("[]")
		database.Seek(0, 0)
	}

	return nil
}

func NewFileSystemStore(database *os.File) (*FileSystemStore, error) {
	err := initializePlayerDBFile(database)
	if err != nil {
		return nil, fmt.Errorf("problem initializing player db file, %v", err)
	}
	league, err := NewLeague(database)
	if err != nil {
		return nil, fmt.Errorf("problem loading player from database %s, %v", database.Name(), err)
	}

	return &FileSystemStore{
		database: json.NewEncoder(&Tape{database}),
		league:   league,
	}, nil
}

func (f *FileSystemStore) GetLeaguePlayers() League {
	sort.Slice(f.league, func(i, j int) bool {
		return f.league[i].Wins > f.league[j].Wins
	})
	return f.league
}

func (f *FileSystemStore) GetPlayerScore(name string) int {
	p := f.league.Find(name)
	if p != nil {
		return p.Wins
	}
	return 0
}

func (f *FileSystemStore) RecordWin(name string) {
	p := f.league.Find(name)
	if p != nil {
		p.Wins++
	} else {
		f.league = append(f.league, player{name, 1})
	}

	f.database.Encode(f.league)
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

type Tape struct {
	file *os.File
}

func (t *Tape) Write(p []byte) (int, error) {
	t.file.Truncate(0)
	t.file.Seek(0, 0)
	return t.file.Write(p)
}
