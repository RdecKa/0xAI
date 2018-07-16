function hexGrid(socket) {
	return new Vue({
		el: '#hex-grid',
		data: {
			socket: socket,
			size: 1,
			grid: []
		},
		methods: {
			initGrid: function (size) {
				this.size = size;
				createHexGrid(this, this.size);
			}
		},
		created: function() {
			initSocket(this, this.socket);
		}
	});
}

/**
 * Adds event listeners to socket in the given Vue instance.
 * @param self   Vue instance
 * @param socket socket to add listeners to
 */
function initSocket(self, socket) {
	socket.addEventListener("open", function (event) {
		console.log("Connection established.");
	});

	socket.addEventListener("message", function (event) {
		let m = event.data;
		console.log("Message from server:", m);
		let response = respondToMessage(self, m);
		if (response !== undefined) {
			socket.send(response);
		}
	});
}

/**
 * Reads the message and returns a response.
 * @param self Vue instance
 * @param msg  string message to be parsed
 */
function respondToMessage(self, msg) {
	let ms = msg.split(" ");
	switch (ms[0]) {
		case "INIT":
			console.log("Init");
			let s = ms[1].split(":");
			self.initGrid(parseInt(s[1]));

			let c = ms[2].split(":")
			switch (c[1]) {
				case "r":
					myColor = colors.RED;
					break;
				case "b":
					myColor = colors.BLUE;
					break;
				default:
					myColor = colors.NONE;
			}
			if (myColor == colors.RED || myColor == colors.BLUE) {
				return "READY";
			} else {
				return "ERROR";
			}
		case "MOVE":
			console.log("Move");
			if (ms[1] !== "<nil>") {
				receiveMove(self, msg.substring(5))
			}
			printGrid(self.grid);
			return nextMove(self.grid, myColor);
		case "END":
			console.log("End");
			receiveMove(self, msg.substring(6));
			printGrid(self.grid);
			myColor = colors.NONE;
			return
		default:
			console.log("Unknown message: '" + msg + "'");
			return
	}
}

/**
 * Decodes a move, changes the grid accordingly.
 * @param self Vue instance
 * @param move string representing a move; example: 'r: (2, 3)'
 */
function receiveMove(self, move) {
	let color;
	switch (move.charAt(0)) {
		case 'r':
			color = colors.RED;
			break;
		case 'b':
			color = colors.BLUE;
			break;
		default:
			console.log("INVALID COLOR '" + move.charAt(0) + "'");
			color = colors.NONE;
	}

	let coords = move.substring(4, move.length - 1).split(", ");

	let x = parseInt(coords[0]);
	let y = parseInt(coords[1]);
	self.grid[y].splice(x, 1, color);
}

/**
 * Creates a new 2D array to be used as a game board.
 * @param self Vue instance
 * @param size size of the grid to be created
 */
function createHexGrid(self, size) {
	self.grid = [];
	for (let i = 0; i < size; i++) {
		self.grid.push([]);
		for (let j = 0; j < size; j++) {
			self.grid[i].push(colors.NONE);
		}
	}
}
