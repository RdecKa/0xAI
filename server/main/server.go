package main

import (
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"regexp"

	"github.com/RdecKa/mcts/game"
	"github.com/RdecKa/mcts/hex"
)

var validPath = regexp.MustCompile("^/(intro|play|sendmove|getmove)/([a-zA-Z0-9]*)$")

var templates = template.Must(template.ParseFiles("server/tmpl/play.html"))

var state hex.State

func makeHandler(fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		a := validPath.FindStringSubmatch(r.URL.Path)
		if a == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r)
	}
}

func playHandler(w http.ResponseWriter, r *http.Request) {
	state = *hex.NewState(byte(6))
	err := templates.ExecuteTemplate(w, "play.html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func getmoveHandler(w http.ResponseWriter, r *http.Request) {
	actions := state.GetPossibleActions()
	a := actions[rand.Intn(len(actions))]
	s, end := makeMove(state, a)
	state = s
	if end {
		w.Write([]byte("END"))
		state = *hex.NewState(byte(state.GetSize()))
	}
	w.Write([]byte(fmt.Sprintf("%v\n", a)))
	w.Write([]byte(fmt.Sprintf("%v", s)))
}

// makeMove returns state in a game after action a has been made in state s and
// a boolean value indicating the end of the game (true if game is finished,
// false otherwise)
func makeMove(s hex.State, a game.Action) (hex.State, bool) {
	ns := s.GetSuccessorState(a).(hex.State)
	if ns.IsGoalState() {
		return ns, true
	}
	return ns, false
}

func main() {
	http.HandleFunc("/play/", makeHandler(playHandler))
	http.HandleFunc("/getmove/", makeHandler(getmoveHandler))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
