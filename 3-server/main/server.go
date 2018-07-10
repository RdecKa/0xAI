package main

import (
	"errors"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"regexp"
	"strconv"

	"github.com/RdecKa/bachleor-thesis/common/game"
	"github.com/RdecKa/bachleor-thesis/common/game/hex"
)

var validPath = regexp.MustCompile("^/((play|sendmove|getmove|static)/([a-zA-Z0-9/.]*))?$")

var templates = template.Must(template.ParseFiles("3-server/tmpl/play.html"))

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

func readCoordinateFromFlag(flag string, r *http.Request) (int, error) {
	s, ok := r.URL.Query()[flag]
	if !ok || len(s[0]) < 1 {
		return -1, errors.New("Coordinate " + flag + " not specified.")
	}
	x, err := strconv.Atoi(s[0])
	if err != nil {
		return -1, errors.New("Coordinate " + flag + " invalid.")
	}
	return x, nil
}

func sendmoveHandler(w http.ResponseWriter, r *http.Request) {
	cs := [2]string{"x", "y"}
	coords := [2]int{}
	for i, f := range cs {
		c, err := readCoordinateFromFlag(f, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		coords[i] = c
	}
	w.Write([]byte(fmt.Sprintf("Place a stone on (%d, %d)", coords[0], coords[1])))
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
	// Register handlers
	http.HandleFunc("/play/", makeHandler(playHandler))
	http.HandleFunc("/getmove/", makeHandler(getmoveHandler))
	http.HandleFunc("/sendmove/", makeHandler(sendmoveHandler))

	// TODO: DELETE / CHANGE
	http.HandleFunc("/", makeHandler(playHandler))

	// Register folder with static content (js, css)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("3-server/static"))))

	// Run server
	log.Println("Server running on loclhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
