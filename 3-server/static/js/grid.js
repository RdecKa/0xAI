const hexSide = 30;
const unitX = 2 * Math.cos(30 * Math.PI / 180) * hexSide;
const unitY = (1 + Math.sin(30 * Math.PI / 180)) * hexSide;
const margin = hexSide * 0.1;

function hexGrid(socket) {
	return new Vue({
		el: '#hex-grid',
		data: {
			socket: socket,
			size: 1,
			grid: []
		},
		computed: {
			boardWidth: function () {
				let width = Math.floor(1.5 * this.size) * unitX + 2 * margin;
				return width + "px";
			},
			boardHeight: function () {
				let height = (this.size - 1) * unitY + 2 * hexSide + 2 * margin;
				return height + "px";
			}
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
 * @param obj    Vue instance
 * @param socket socket to add listeners to
 */
function initSocket(obj, socket) {
	socket.addEventListener("open", function (event) {
		console.log("Connection established.");
	});

	socket.addEventListener("message", function (event) {
		let m = event.data;
		console.log("Message from server:", m);
		let response = respondToMessage(obj, m);
		if (response !== undefined) {
			socket.send(response);
		}
	});
}

/**
 * Reads the message and returns a response.
 * @param obj Vue instance
 * @param msg string message to be parsed
 */
function respondToMessage(obj, msg) {
	let ms = msg.split(" ");
	switch (ms[0]) {
		case "INIT":
			console.log("Init");
			let s = ms[1].split(":");
			obj.initGrid(parseInt(s[1]));

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
				receiveMove(obj, msg.substring(5))
			}
			printGrid(obj.grid);
			return nextMove(obj.grid, myColor);
		case "END":
			console.log("End");
			receiveMove(obj, msg.substring(6));
			printGrid(obj.grid);
			myColor = colors.NONE;
			return
		default:
			console.log("Unknown message: '" + msg + "'");
			return
	}
}

/**
 * Decodes a move, changes the grid accordingly.
 * @param obj  Vue instance
 * @param move string representing a move; example: 'r: (2, 3)'
 */
function receiveMove(obj, move) {
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
	obj.grid[y].splice(x, 1, color);
}

/**
 * Creates a new 2D array to be used as a game board.
 * @param obj  Vue instance
 * @param size size of the grid to be created
 */
function createHexGrid(obj, size) {
	obj.grid = [];
	for (let i = 0; i < size; i++) {
		obj.grid.push([]);
		for (let j = 0; j < size; j++) {
			obj.grid[i].push(colors.NONE);
		}
	}
}
