package main

import (
	"html/template"
	"log"
	"math"
	"net/http"
	"regexp"
	"time"

	"github.com/RdecKa/bachleor-thesis/3-server/hexgame"
	"github.com/RdecKa/bachleor-thesis/3-server/hexplayer"
	"github.com/RdecKa/bachleor-thesis/common/game/hex"
)

var validPath = regexp.MustCompile("^/((play|select|static|ws)/([a-zA-Z0-9/.]*))?$")

var templates = template.Must(template.New("").Delims("[[", "]]").ParseFiles(
	"3-server/tmpl/play.html",
	"3-server/tmpl/select.html"))

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

func selectHandler(w http.ResponseWriter, r *http.Request) {
	err := templates.ExecuteTemplate(w, "select.html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	args := r.URL.Query()
	red, okRed := args["red"]
	blue, okBlue := args["blue"]
	pair := [2]hexplayer.HexPlayer{} // 0 - red, 1 - blue

	var rFunc, bFunc func(w http.ResponseWriter, r *http.Request, color hex.Color) hexplayer.HexPlayer

	if okRed && red[0] == "human" {
		rFunc = createHumanPlayer
	} else if okRed && red[0] == "mcts" {
		rFunc = createMCTSplayer
	}

	if okBlue && blue[0] == "human" && red[0] != "human" {
		bFunc = createHumanPlayer
	} else if okBlue && blue[0] == "mcts" && red[0] != "mcts" {
		bFunc = createMCTSplayer
	}

	if rFunc == nil || bFunc == nil {
		log.Println("Wrong or missing arguments for players. Using default.")
		pair[0] = createHumanPlayer(w, r, hex.Red)
		pair[1] = createMCTSplayer(w, r, hex.Blue)
	} else {
		pair[0] = rFunc(w, r, hex.Red)
		pair[1] = bFunc(w, r, hex.Blue)
	}

	go hexgame.Play(pair, 10)
}

func createHumanPlayer(w http.ResponseWriter, r *http.Request, color hex.Color) hexplayer.HexPlayer {
	conn, err := hexplayer.OpenConn(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err)
		return nil
	}
	return hexplayer.CreateHumanPlayer(conn, color)
}

func createMCTSplayer(w http.ResponseWriter, r *http.Request, color hex.Color) hexplayer.HexPlayer {
	return hexplayer.CreateMCTSplayer(color, math.Sqrt(2), time.Duration(2)*time.Second, 10)
}

func main() {
	// Register handlers
	http.HandleFunc("/play/", makeHandler(playHandler))
	http.HandleFunc("/select/", makeHandler(selectHandler))

	http.HandleFunc("/ws/", makeHandler(wsHandler))

	// TODO: DELETE / CHANGE
	http.HandleFunc("/", makeHandler(playHandler))

	// Register folder with static content (js, css)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("3-server/static"))))

	// Run server
	log.Println("Server running on", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
