package mcts

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/RdecKa/bachleor-thesis/common/game/hex"
	"github.com/RdecKa/bachleor-thesis/common/tree"
)

// workerChan has a list of all chanells that workers use
type workerChan struct {
	assign     chan *MCTS        // Sending new tasks to workers
	gather     chan []*tree.Node // Receiving results from workers
	e          chan error        // Receiving errors from workers
	quit       chan struct{}     // Signaling workers to stop
	terminated chan struct{}     // Signaling that worker terminated
}

// RunMCTSinParallel takes care of running MCTS in parallel. It creates
// numWorkers workers - each of them runs one instance of MCTS at once.
// Iterations of MCTS are run on board of size boardSize for timeToRun. mc is
// the initialised search that is completed first.
// If gameLengthImportant is true, then a goal state with a shorter path to
// victory gets a higher estimated value than a goal state with a longer path.
func RunMCTSinParallel(numWorkers, boardSize int, treasholdN uint, timeToRun time.Duration,
	outputFolder, patFileName string, mc *MCTS, gameLengthImportant bool) {
	var err error

	assign := make(chan *MCTS, numWorkers)
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

	filePrefix := fmt.Sprintf("sample_%02d_%d", boardSize, int(timeToRun.Seconds()))
	fileNameNoEnding := fmt.Sprintf("%s%s", outputFolder, filePrefix)

	// Create a log file
	logFile, err := os.Create(fileNameNoEnding + ".log")
	if err != nil {
		panic(err)
	}

	// Create workers
	for w := 0; w < numWorkers; w++ {
		// Create a file for a worker to write learning samples
		fileName := fmt.Sprintf("%s_%d.in", fileNameNoEnding, w)
		f, err := os.Create(fileName)
		if err != nil {
			panic(err)
		}
		defer f.Close()

		// Create a file for a worker to store details about the search
		fDet, err := os.Create(fileName + ".det")
		if err != nil {
			panic(err)
		}
		defer fDet.Close()

		// Start a worker process
		go worker(w, timeToRun, boardSize, treasholdN, f, fDet, logFile,
			patFileName, &wc, gameLengthImportant)
	}

	candidateList := NewCandidateList(boardSize)

	numExpansions := 1
	finished, quitted := false, false
	tasksAssigned, tasksFinished := 1, 0
	for !finished {
		logFile.WriteString(fmt.Sprintf("Queue:\n%v", candidateList))
		select {
		case <-kill:
			quitted = true
		case newCandidates := <-gather:
			// Get a returned value from a worker
			tasksFinished++
			candidateList.AddCandidates(newCandidates)
			numExpansions++

			if !quitted {
				// Add as many new tasks as there are free spots in the assign
				// channel (if there are at least that many tasks)
				for t := len(assign); t < numWorkers; t++ {
					newTask := candidateList.GetNextCandidateToExpand()
					if newTask == nil {
						break
					}
					tasksAssigned++
					assign <- mc.ContinueMCTSFromNode(newTask)
				}
			}

			if tasksAssigned == tasksFinished && (candidateList.IsEmpty() || quitted) {
				logFile.WriteString(fmt.Sprintln("All tasks finished (or QUIT signal received)!"))
				for w := 0; w < numWorkers; w++ {
					quit <- struct{}{}
				}
				finished = true
			}
		case err = <-e:
			panic(err)
		}
	}

	// Wait for all workers to terminate
	for w := 0; w < numWorkers; w++ {
		<-terminated
	}
}

// worker waits for tasks and executes them in an infinite loop until the quit
// signal
func worker(id int, timeToRun time.Duration, boardSize int, treasholdN uint,
	outputFile, outputFileDet, logFile *os.File, patFileName string,
	wc *workerChan, gameLengthImportant bool) {

	var mc *MCTS
	taskID := 0
	gridChan, stopChan, resultChan := hex.CreatePatChecker(patFileName)
	outputFile.WriteString(hex.GetHeaderCSV())
	for {
		select {
		case mc = <-wc.assign:
			outputFile.WriteString(fmt.Sprintf("# Search ID %d\n", taskID))
			outputFileDet.WriteString(fmt.Sprintf("# Search ID %d started from:\n%v\n", taskID, mc.GetInitialNode()))
			expCand, err := RunMCTS(mc, id, timeToRun, boardSize, treasholdN,
				outputFile, logFile, gridChan, resultChan, gameLengthImportant)
			if err != nil {
				wc.e <- err
			}
			logFile.WriteString(fmt.Sprintf("Worker %d finished task %d\n", id, taskID))
			wc.gather <- expCand
			taskID++
		case <-wc.quit:
			logFile.WriteString(fmt.Sprintf("Worker %d terminated\n", id))
			stopChan <- struct{}{} // Stop the pattern checker
			wc.terminated <- struct{}{}
			return
		}
	}
}

// sampleArrayOfNodes accepts an array of elements. It returns a new array that
// contains some elements from the old array. Each element of the old array is
// copied to the new array with a probability p, otherwise it is discarded.
func sampleArrayOfNodes(oldArr []*tree.Node, p float64) []*tree.Node {
	newArr := make([]*tree.Node, 0)
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
	fmt.Println("Enter 'q' to quit")
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
