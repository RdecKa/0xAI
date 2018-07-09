package mcts

import (
	"fmt"
	"math/rand"
	"os"

	"github.com/RdecKa/common/game/hex"
	"github.com/RdecKa/common/tree"
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
// numIterations of MCTS are run on board of size boardSize. mc is the
// initialised search that is completed first.
func RunMCTSinParallel(numWorkers, boardSize, numIterations int, outputFolder string, mc *MCTS) {
	var expCand []*tree.Node // Array of possible candidates for continuing MCTS
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

	filePrefix := fmt.Sprintf("sample_%02d_%d", boardSize, numIterations)
	fileNameNoEnding := fmt.Sprintf("%s%s", outputFolder, filePrefix)

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
		go worker(w, numIterations, boardSize, f, fDet, logFile, &wc)
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
			newCandidates = sampleArrayOfNodes(newCandidates, 1/float64(boardSize*boardSize))
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
}

// worker waits for tasks and executes them in an infinite loop
func worker(id, numIterations, boardSize int, outputFile, outputFileDet, logFile *os.File, wc *workerChan) {
	var mc *MCTS
	taskID := 0
	gridChan, stopChan, resultChan := hex.CreatePatChecker()
	outputFile.WriteString(hex.GetHeaderCSV())
	for {
		logFile.WriteString(fmt.Sprintf("Worker %d waiting for a task\n", id))
		select {
		case mc = <-wc.assign:
			logFile.WriteString(fmt.Sprintf("Worker %d executing task\n", id))
			outputFile.WriteString(fmt.Sprintf("# Search ID %d\n", taskID))
			outputFileDet.WriteString(fmt.Sprintf("# Search ID %d started from:\n%v\n", taskID, mc.GetInitialNode()))
			expCand, err := RunMCTS(mc, id, numIterations, boardSize, outputFile, logFile, gridChan, resultChan)
			if err != nil {
				wc.e <- err
			}
			logFile.WriteString(fmt.Sprintf("Worker %d finished a task\n", id))
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
