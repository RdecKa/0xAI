package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"math"
	"net/http"
	_ "net/http/pprof"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/RdecKa/bachleor-thesis/common/game/hex"
	"github.com/RdecKa/bachleor-thesis/server/cmpr"
	"github.com/RdecKa/bachleor-thesis/server/hexgame"
	"github.com/RdecKa/bachleor-thesis/server/hexplayer"
	"github.com/gorilla/websocket"
)

var validPath = regexp.MustCompile("^/((play|select|static|ws)/([a-zA-Z0-9/.]*))?$")

var templates = template.Must(template.New("").Delims("[[", "]]").ParseFiles(
	"server/tmpl/play.html",
	"server/tmpl/select.html"))

const addr = "localhost:8080"
const patternFile = "common/game/hex/patterns.txt"

var resultsDir = "data/cmpr/"

const defaultBoardSize = 7
const defaultNumGames = 1
const defaultTime = 1

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
	boardSizeString, okBoardSize := args["size"]
	numGamesString, okNumGames := args["numgames"]
	redTimeString, okRedTime := args["redtime"]
	blueTimeString, okBlueTime := args["bluetime"]

	wa := okWatch && watch[0] == "false"
	pair := [2]hexplayer.HexPlayer{} // 0 - red, 1 - blue

	var boardSize, numGames int
	var err error
	if !okBoardSize {
		boardSize = defaultBoardSize
	} else {
		boardSize, err = strconv.Atoi(boardSizeString[0])
		if err != nil {
			log.Println(err)
			boardSize = defaultBoardSize
		}
	}
	if !okNumGames {
		numGames = defaultNumGames
	} else {
		numGames, err = strconv.Atoi(numGamesString[0])
		if err != nil {
			log.Println(err)
			numGames = defaultNumGames
		}
	}

	var redTime, blueTime int
	if !okRedTime {
		redTime = defaultTime
	} else {
		redTime, err = strconv.Atoi(redTimeString[0])
		if err != nil {
			log.Println(err)
			redTime = defaultTime
		}
	}
	if !okBlueTime {
		blueTime = defaultTime
	} else {
		blueTime, err = strconv.Atoi(blueTimeString[0])
		if err != nil {
			log.Println(err)
			blueTime = defaultTime
		}
	}

	conn, err := hexplayer.OpenConn(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err)
		return
	}

	var rFunc, bFunc func(hex.Color, *websocket.Conn, int, bool, hexplayer.PlayerType) hexplayer.HexPlayer

	if okRed && red[0] == "human" {
		rFunc = createHumanPlayer
	} else if okRed && red[0] == "mcts" {
		rFunc = createMCTSplayer
	} else if okRed && red[0][:2] == "ab" {
		rFunc = createAbPlayer
	}

	if okBlue && blue[0] == "human" && red[0] != "human" {
		bFunc = createHumanPlayer
	} else if okBlue && blue[0] == "mcts" {
		bFunc = createMCTSplayer
	} else if okBlue && blue[0][:2] == "ab" {
		bFunc = createAbPlayer
	}

	if rFunc == nil || bFunc == nil {
		log.Println("Wrong or missing arguments for players. Using default.")
		wa = false
		pair[0] = createHumanPlayer(hex.Red, conn, redTime, wa, hexplayer.GetPlayerTypeFromString(red[0]))
		pair[1] = createMCTSplayer(hex.Blue, conn, blueTime, wa, hexplayer.GetPlayerTypeFromString(blue[0]))
	} else {
		pair[0] = rFunc(hex.Red, conn, redTime, wa, hexplayer.GetPlayerTypeFromString(red[0]))
		pair[1] = bFunc(hex.Blue, conn, blueTime, wa, hexplayer.GetPlayerTypeFromString(blue[0]))
	}

	c := conn
	if wa {
		c = nil
	}

	go hexgame.Play(boardSize, pair, numGames, c, nil, nil, resultsDir)
}

func createHumanPlayer(color hex.Color, conn *websocket.Conn, _ int, _ bool, _ hexplayer.PlayerType) hexplayer.HexPlayer {
	return hexplayer.CreateHumanPlayer(conn, color)
}

func createMCTSplayer(color hex.Color, _ *websocket.Conn, secondsPerAction int, allowResignation bool, _ hexplayer.PlayerType) hexplayer.HexPlayer {
	return hexplayer.CreateMCTSplayer(color, math.Sqrt(2), time.Duration(secondsPerAction)*time.Second, 10, allowResignation)
}

func createAbPlayer(color hex.Color, conn *websocket.Conn, secondsPerAction int, allowResignation bool, subtype hexplayer.PlayerType) hexplayer.HexPlayer {
	return hexplayer.CreateAbPlayer(color, conn, time.Duration(secondsPerAction)*time.Second, allowResignation, patternFile, true, subtype)
}

func comparePlayers() {
	matches := []cmpr.MatchSetup{
		cmpr.CreateMatch(7, 2, hexplayer.MctsType, hexplayer.AbDtType, 1, 1, patternFile),
		cmpr.CreateMatch(7, 2, hexplayer.MctsType, hexplayer.AbLrType, 1, 1, patternFile),
		cmpr.CreateMatch(7, 2, hexplayer.AbDtType, hexplayer.AbLrType, 1, 1, patternFile),
	}

	cmpr.RunAll(matches, resultsDir)
}

func main() {
	t := time.Now()
	resultsDir += t.Format("20060102T150405") + "/"
	err := os.MkdirAll(resultsDir, os.ModePerm)
	if err != nil {
		fmt.Println(err)
	}

	pOnlyCompare := flag.Bool("cmpr", false, "Run test matches between players")
	flag.Parse()
	if *pOnlyCompare {
		fmt.Println("Running comparisons")
		comparePlayers()
		return
	}

	// Register handlers
	http.HandleFunc("/play/", makeHandler(playHandler))
	http.HandleFunc("/select/", makeHandler(selectHandler))
	http.HandleFunc("/ws/", makeHandler(wsHandler))

	// TODO: DELETE / CHANGE
	//http.HandleFunc("/", makeHandler(playHandler))

	// Register folder with static content (js, css)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("server/static"))))

	// Run server
	log.Println("Server running on", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
