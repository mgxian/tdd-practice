package players

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestGETPlayers(t *testing.T) {
	store := &stubPlayerStore{
		map[string]int{
			"Pepper": 20,
			"Floyd":  10,
		},
		nil,
		nil,
	}
	server := newPlayerServer(store)
	t.Run("Return Pepper's score", func(t *testing.T) {
		request := newGetPlayerScore("Pepper")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusOK)
		assertResponseBody(t, response.Body.String(), "20")
	})

	t.Run("Return Floyd's score", func(t *testing.T) {
		request := newGetPlayerScore("Floyd")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusOK)
		assertResponseBody(t, response.Body.String(), "10")

	})

	t.Run("Return 404 if player doesn't exist", func(t *testing.T) {
		request := newGetPlayerScore("unknown")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusNotFound)
	})
}

func TestStoreWins(t *testing.T) {
	store := &stubPlayerStore{
		map[string]int{},
		nil,
		nil,
	}
	server := newPlayerServer(store)

	t.Run("Record wins on POST", func(t *testing.T) {
		player := "Pepper"
		request := newPostWinRequest(player)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)
		assertStatus(t, response.Code, http.StatusAccepted)

		if len(store.winCalls) != 1 {
			t.Errorf("got %d calls to RecordWin want %d", len(store.winCalls), 1)
		}

		if store.winCalls[0] != player {
			t.Errorf("did not record correct winer got %s want %s", store.winCalls, player)
		}
	})
}

func TestRecordingWindsAndRetrievingThem(t *testing.T) {
	store := newInMemoryPlayerStore()
	server := newPlayerServer(store)
	player := "Pepper"

	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))
	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))
	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))

	response := httptest.NewRecorder()
	server.ServeHTTP(response, newGetPlayerScore(player))
	assertStatus(t, response.Code, http.StatusOK)
	assertResponseBody(t, response.Body.String(), "3")
}

func TestLeague(t *testing.T) {
	t.Run("Return league as JSON", func(t *testing.T) {
		wantedLeague := []player{
			{"Cleo", 32},
			{"Chris", 20},
			{"Tiest", 14},
		}
		store := &stubPlayerStore{nil, nil, wantedLeague}
		server := newPlayerServer(store)

		request, _ := http.NewRequest(http.MethodGet, "/league", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)
		assertStatus(t, response.Code, http.StatusOK)

		gotLeague := getLeagueFromResponse(t, response.Body)
		assertLeague(t, gotLeague, wantedLeague)
		assertContentType(t, response.Result().Header.Get("Content-Type"), jsonContentType)
	})
}

func getLeagueFromResponse(t *testing.T, body io.Reader) []player {
	t.Helper()
	var league []player
	err := json.NewDecoder(body).Decode(&league)
	if err != nil {
		t.Fatalf("Unable to decode %q, %v", body, err)
	}
	return league
}

func assertContentType(t *testing.T, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("response did not have content-type of %q, got %q", got, want)
	}
}

func assertLeague(t *testing.T, got, want []player) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func newGetPlayerScore(name string) *http.Request {
	request, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/players/%s", name), nil)
	return request
}

func newPostWinRequest(name string) *http.Request {
	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/players/%s", name), nil)
	return req
}

func assertResponseBody(t *testing.T, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func assertStatus(t *testing.T, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("got status %d, want %d", got, want)
	}
}
