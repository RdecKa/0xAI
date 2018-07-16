package main

import (
	"html/template"
	"log"
	"math"
	"net/http"
	"regexp"

	"github.com/RdecKa/bachleor-thesis/3-server/hexgame"
	"github.com/RdecKa/bachleor-thesis/3-server/hexplayer"
	"github.com/RdecKa/bachleor-thesis/common/game/hex"
)

var validPath = regexp.MustCompile("^/((play|sendmove|getmove|static|ws)/([a-zA-Z0-9/.]*))?$")

var templates = template.Must(template.New("").Delims("[[", "]]").ParseFiles("3-server/tmpl/play.html"))

const addr = "localhost:8080"

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
	err := templates.ExecuteTemplate(w, "play.html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	pair := [2]hexplayer.HexPlayer{}
	colors := [2]hex.Color{hex.Red, hex.Blue}
	// In Human:Human version players share the same websocket
	conn, err := hexplayer.OpenConn(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	/*for i := 0; i < 2; i++ {
		hp := hexplayer.CreateHumanPlayer(conn, colors[i])
		pair[i] = hp
	}*/
	pair[0] = hexplayer.CreateHumanPlayer(conn, colors[0])
	pair[1] = hexplayer.CreateMCTSplayer(colors[1], math.Sqrt(2), 5000, 10)

	go hexgame.Play(pair, 1)
}

func main() {
	// Register handlers
	http.HandleFunc("/play/", makeHandler(playHandler))

	http.HandleFunc("/ws/", makeHandler(wsHandler))

	// TODO: DELETE / CHANGE
	http.HandleFunc("/", makeHandler(playHandler))

	// Register folder with static content (js, css)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("3-server/static"))))

	// Run server
	log.Println("Server running on", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
