package main

import (
	"flag"
	"fmt"
	"math"
	"math/rand"
	"os"

	"github.com/RdecKa/mcts/hex"
	"github.com/RdecKa/mcts/mcts"
	"github.com/RdecKa/mcts/tree"
)

type workerChan struct {
	assign     chan *mcts.MCTS
	gather     chan []*tree.Node
	e          chan error
	done       chan struct{}
	terminated chan struct{}
}

func main() {
	boardSize := flag.Int("size", 3, "Board size")
	numIterations := flag.Int("iter", 10000, "Number of iterations")
	indentJSON := flag.Bool("indent", false, "Indent JSON output")
	output := flag.String("output", ".", "Output folder")
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
	done := make(chan struct{}, numWorkers)
	terminated := make(chan struct{}, numWorkers)
	defer close(assign)
	defer close(gather)
	defer close(e)
	defer close(done)
	defer close(terminated)

	wc := workerChan{assign, gather, e, done, terminated}

	assign <- mc // Send first task

	// Create workers
	for w := 0; w < numWorkers; w++ {
		go worker(w, *numIterations, *boardSize, *output, &wc)
	}

	finished := false
	tasksAssigned, tasksFinished := 1, 0
	for !finished {
		fmt.Printf("Queue len: %d\n", len(expCand))
		select {
		case newCandidates := <-gather:
			// Get a returned value from a worker
			tasksFinished++
			newCandidates = sampleArrayOfNodes(newCandidates, 1/float64((*boardSize)*(*boardSize)))
			expCand = append(expCand, newCandidates...)

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

			if len(expCand) == 0 && tasksAssigned == tasksFinished {
				fmt.Println("All tasks finished!")
				for w := 0; w < numWorkers; w++ {
					done <- struct{}{}
				}
				finished = true
				break
			}
		case err = <-e:
			fmt.Println(err)
			os.Exit(1)
		}
	}

	// Write JSON
	filePrefix := fmt.Sprintf("out_%02d_%d", *boardSize, *numIterations)
	err = mcts.WriteToFile(*root, *output, filePrefix, *indentJSON)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Wait for all workers to terminate
	for w := 0; w < numWorkers; w++ {
		<-terminated
	}
}

func worker(id, numIterations, boardSize int, output string, wc *workerChan) {
	var mc *mcts.MCTS
	for {
		fmt.Printf("Worker %d waiting for a task\n", id)
		select {
		case mc = <-wc.assign:
		case <-wc.done:
			fmt.Printf("Worker %d terminated\n", id)
			wc.terminated <- struct{}{}
			return
		}
		fmt.Printf("Worker %d executing task\n", id)
		expCand, err := runMCTS(mc, numIterations, boardSize, output)
		if err != nil {
			wc.e <- err
		}
		fmt.Printf("Worker %d finished a task\n", id)
		wc.gather <- expCand
	}
}

// runMCTS executes numIterations iterations of MCTS, given initialised MCTS
func runMCTS(mc *mcts.MCTS, numIterations, boardSize int, output string) ([]*tree.Node, error) {
		fmt.Printf("Starting new MCTS from node\n%v\n", mc.GetInitialNode())
	for i := 0; i < numIterations; i++ {
		if i > 0 && i%10000 == 0 {
			fmt.Printf("Finished iteration %d\n", i)
		}
		mc.RunIteration()
	}

		// Write input-output pairs for supervised machine learning, generate
		// new nodes tocontinue MCTS
	filePrefix := fmt.Sprintf("sample_%02d_%d", boardSize, numIterations)
	expCand, err := mc.GenSamples(output, filePrefix, 100)
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
