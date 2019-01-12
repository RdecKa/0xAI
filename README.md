# 0xAI

Code for my bachelor thesis **Using Monte Carlo tree search and machine learning to learn a heuristic function**.

The thesis with the description of the project (in Slovene) will be avaliable soon.

## How to run the code

The easiest way is to use the provided *Makefile*.

1. If you want, you can change some values in *Makefile*, for example:
  * SIZE - the size of the Hex grid
  * TIME - how much time can a single MCTS run (in seconds)
  * WORKERS - how many goroutines should be created to run MCTS in parallel
  * THRESHOLD_N - how many times should a node of MCTS tree be visited to be used as a learning sample.
1. Run `make` in the root directory.
1. Monte Carlo tree search (MCTS) will start generating learning samples in folder *data/SIZE/mcts/run-START_TIME/*. When you are satisfied with the number of samples, type `q` and press Enter.
1. Machine learning (ML) phase will start and generate code with evaluation functions. Just wait.
1. Code for the server will be compiled. Follow the instructions on the screen to run it.
1. Open *localhost:8080/select/* and play! Or just watch.


It is possible to run only specific parts of the project.

* `make mcts` will run only MCTS phase.
* `make ml START_TIME=TIME` will only run ML phase using learning samples from *data/SIZE/mcts/run-TIME/*.
* `make serv` will compile the server with the heuristic functions from the last run of ML phase.
