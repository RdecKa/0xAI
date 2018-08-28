package main

import _ "net/http/pprof"
import (
	"html/template"
	"log"
	"math"
	"net/http"
	"regexp"
	"time"

	"github.com/RdecKa/bachleor-thesis/common/game/hex"
	"github.com/RdecKa/bachleor-thesis/server/hexgame"
	"github.com/RdecKa/bachleor-thesis/server/hexplayer"
	"github.com/gorilla/websocket"
)

var validPath = regexp.MustCompile("^/((play|select|static|ws)/([a-zA-Z0-9/.]*))?$")

var templates = template.Must(template.New("").Delims("[[", "]]").ParseFiles(
	"server/tmpl/play.html",
	"server/tmpl/select.html"))

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
	watch, okWatch := args["watch"]
	wa := okWatch && watch[0] == "false"
	pair := [2]hexplayer.HexPlayer{} // 0 - red, 1 - blue

	conn, err := hexplayer.OpenConn(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err)
		return
	}

	var rFunc, bFunc func(hex.Color, *websocket.Conn, bool) hexplayer.HexPlayer

	if okRed && red[0] == "human" {
		rFunc = createHumanPlayer
	} else if okRed && red[0] == "mcts" {
		rFunc = createMCTSplayer
	} else if okRed && red[0] == "ab" {
		rFunc = createAbPlayer
	}

	if okBlue && blue[0] == "human" && red[0] != "human" {
		bFunc = createHumanPlayer
	} else if okBlue && blue[0] == "mcts" {
		bFunc = createMCTSplayer
	} else if okBlue && blue[0] == "ab" {
		bFunc = createAbPlayer
	}

	if rFunc == nil || bFunc == nil {
		log.Println("Wrong or missing arguments for players. Using default.")
		wa = false
		pair[0] = createHumanPlayer(hex.Red, conn, wa)
		pair[1] = createMCTSplayer(hex.Blue, conn, wa)
	} else {
		pair[0] = rFunc(hex.Red, conn, wa)
		pair[1] = bFunc(hex.Blue, conn, wa)
	}

	c := conn
	if wa {
		c = nil
	}

	go hexgame.Play(pair, 15, c)
}

func createHumanPlayer(color hex.Color, conn *websocket.Conn, allowResignation bool) hexplayer.HexPlayer {
	return hexplayer.CreateHumanPlayer(conn, color)
}

func createMCTSplayer(color hex.Color, conn *websocket.Conn, allowResignation bool) hexplayer.HexPlayer {
	return hexplayer.CreateMCTSplayer(color, math.Sqrt(2), time.Duration(1)*time.Second, 10, allowResignation)
}

func createAbPlayer(color hex.Color, conn *websocket.Conn, allowResignation bool) hexplayer.HexPlayer {
	return hexplayer.CreateAbPlayer(color, conn, time.Duration(1)*time.Second, allowResignation, "common/game/hex/patterns.txt", false)
}

func main() {
	// Register handlers
	http.HandleFunc("/play/", makeHandler(playHandler))
	http.HandleFunc("/select/", makeHandler(selectHandler))

	http.HandleFunc("/ws/", makeHandler(wsHandler))

	// TODO: DELETE / CHANGE
	http.HandleFunc("/", makeHandler(playHandler))

	// Register folder with static content (js, css)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("server/static"))))

	// Run server
	log.Println("Server running on", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
