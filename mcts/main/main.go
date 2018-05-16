package main

import (
	"flag"
	"fmt"
	"math"
	"math/rand"
	"os"
	"time"

	"github.com/RdecKa/mcts/hex"
	"github.com/RdecKa/mcts/mcts"
	"github.com/RdecKa/mcts/tree"
)

type workerChan struct {
	assign     chan *mcts.MCTS   // Sending new tasks to workers
	gather     chan []*tree.Node // Receiving results from workers
	e          chan error        // Receiving errors from workers
	quit       chan struct{}     // Signaling workers to stop
	terminated chan struct{}     // Signaling that worker terminated
}

func main() {
	boardSize := flag.Int("size", 3, "Board size")
	numIterations := flag.Int("iter", 10000, "Number of iterations")
	indentJSON := flag.Bool("indent", false, "Indent JSON output")
	outputFolder := flag.String("output", ".", "Output folder")
	flag.Parse()
	fmt.Printf("Using boardSize = %d, numIterations = %d\n", *boardSize, *numIterations)

	initState := hex.NewState(byte(*boardSize))
	explorationFactor := math.Sqrt(2)
	minBeforeExpand := uint(10)
	mc := mcts.InitMCTS(*initState, explorationFactor, minBeforeExpand)
	root := mc

	var expCand []*tree.Node // Array of possible candidates for continuing MCTS
	var err error

	numWorkers := 2 // Number of goroutines to run in parallel

	assign := make(chan *mcts.MCTS, numWorkers)
	gather := make(chan []*tree.Node, numWorkers)
	e := make(chan error, 1)
	quit := make(chan struct{}, numWorkers)
	terminated := make(chan struct{}, numWorkers)
	kill := make(chan struct{}, 1)
	defer close(assign)
	defer close(gather)
	defer close(e)
	defer close(quit)
	defer close(terminated)
	defer close(kill)

	wc := workerChan{assign, gather, e, quit, terminated}

	assign <- mc // Send first task

	// Create a boss
	go boss(kill)

	t := time.Now()
	filePrefix := fmt.Sprintf("sample_%02d_%d_%s", *boardSize, *numIterations, t.Format("20060102T150405"))
	fileNameNoEnding := fmt.Sprintf("%s%s", *outputFolder, filePrefix)

	// Create a log file
	logFile, err := os.Create(fileNameNoEnding + ".log")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Create workers
	for w := 0; w < numWorkers; w++ {
		// Create a file for a worker to write learning samples
		fileName := fmt.Sprintf("%s_%d.in", fileNameNoEnding, w)
		f, err := os.Create(fileName)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
	}
		defer f.Close()

		// Create a file for a worker to store details about the search
		fDet, err := os.Create(fileName + ".det")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		defer fDet.Close()

		// Start a worker process
		go worker(w, *numIterations, *boardSize, f, fDet, logFile, &wc)
	}

	finished, quitted := false, false
	tasksAssigned, tasksFinished := 1, 0
	for !finished {
		logFile.WriteString(fmt.Sprintf("Queue len: %d\n", len(expCand)))
		select {
		case <-kill:
			quitted = true
		case newCandidates := <-gather:
			// Get a returned value from a worker
			tasksFinished++
			newCandidates = sampleArrayOfNodes(newCandidates, 1/float64((*boardSize)*(*boardSize)))
			expCand = append(expCand, newCandidates...)

			if !quitted {
			// Add as many new tasks as there are free spots in the assign
			// channel (if there are at least that many tasks)
			for t := len(assign); t < numWorkers; t++ {
				if len(expCand) <= 0 {
					break
				}
				tasksAssigned++
				newTask := rand.Intn(len(expCand))
				assign <- mc.ContinueMCTSFromNode(expCand[newTask])
				expCand = append(expCand[:newTask], expCand[newTask+1:]...)
			}
			}

			if (len(expCand) == 0 || quitted) && tasksAssigned == tasksFinished {
				logFile.WriteString(fmt.Sprintln("All tasks finished (or QUIT signal received)!"))
				for w := 0; w < numWorkers; w++ {
					quit <- struct{}{}
				}
				finished = true
			}
		case err = <-e:
			fmt.Println(err)
			os.Exit(1)
		}
	}

	// Wait for all workers to terminate
	for w := 0; w < numWorkers; w++ {
		<-terminated
	}

	// Write JSON
	filePrefix = fmt.Sprintf("out_%02d_%d", *boardSize, *numIterations)
	err = mcts.WriteToFile(*root, *outputFolder, filePrefix, *indentJSON)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
}
}

// worker waits for tasks and executes them in an infinite loop
func worker(id, numIterations, boardSize int, outputFile, outputFileDet, logFile *os.File, wc *workerChan) {
	var mc *mcts.MCTS
	taskID := 0
	for {
		logFile.WriteString(fmt.Sprintf("Worker %d waiting for a task\n", id))
		select {
		case mc = <-wc.assign:
			logFile.WriteString(fmt.Sprintf("Worker %d executing task\n", id))
			outputFile.WriteString(fmt.Sprintf("# Search ID %d\n", taskID))
			outputFileDet.WriteString(fmt.Sprintf("# Search ID %d started from:\n%v\n", taskID, mc.GetInitialNode()))
			expCand, err := runMCTS(mc, id, numIterations, boardSize, outputFile, logFile)
			if err != nil {
				wc.e <- err
			}
			logFile.WriteString(fmt.Sprintf("Worker %d finished a task\n", id))
			wc.gather <- expCand
			taskID++
		case <-wc.quit:
			logFile.WriteString(fmt.Sprintf("Worker %d terminated\n", id))
			wc.terminated <- struct{}{}
			return
		}
	}
}

// runMCTS executes numIterations iterations of MCTS, given initialised MCTS
func runMCTS(mc *mcts.MCTS, workerID, numIterations, boardSize int, outputFile, logFile *os.File) ([]*tree.Node, error) {
	for i := 0; i < numIterations; i++ {
		if i > 0 && i%10000 == 0 {
			logFile.WriteString(fmt.Sprintf("Worker %d finished iteration %d\n", workerID, i))
		}
		mc.RunIteration()
	}

	// Write input-output pairs for supervised machine learning, generate
	// new nodes to continue MCTS
	expCand, err := mc.GenSamples(outputFile, 100)
	if err != nil {
		return nil, err
	}
	return expCand, nil
}

// sampleArrayOfNodes accepts an array of elements. It returns a new array that
// contains some elements from the old array. Each element of the old array is
// copied to the new array with a probability p, otherwise it is discarded.
func sampleArrayOfNodes(oldArr []*tree.Node, p float64) []*tree.Node {
	newArr := make([]*tree.Node, 0, int(p*float64(len(oldArr))))
	for _, el := range oldArr {
		if rand.Float64() < p {
			newArr = append(newArr, el)
		}
	}
	return newArr
}

// boss waits for a command from a console. If quit signal is received, it stops
// the search (sends a kill signal)
func boss(kill chan struct{}) {
	var input string
	for {
		fmt.Scanln(&input)
		switch input {
		case "q":
			fallthrough
		case "quit":
			fmt.Printf("QUIT signal received. Quitting ...\n")
			kill <- struct{}{}
			return
		default:
			fmt.Printf("Unknown command '%s'\n", input)
		}
	}
}
