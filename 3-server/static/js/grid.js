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
			},
			getPoints: function (colIndex, rowIndex) {
				return getPointsForPolygon(rowIndex, colIndex);
			},
			getClass: function (colIndex, rowIndex) {
				switch (this.grid[rowIndex][colIndex]) {
					case "r":
						return "cell red";
					case "b":
						return "cell blue";
					case ".":
						return "cell empty";
				}
			},
			getBoardWidth: function () {
				return this.size * 150;
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

/**
 * Returns an array of two-ary arrays (coordinates) that represent vertices of a
 * hexagon.
 * @param rowIndex index of a row that contains the cell
 * @param colIndex index of a column that contains the cell
 */
function getPointsForPolygon(rowIndex, colIndex) {
	let centerX = (colIndex  + rowIndex / 2) * unitX;
	let centerY = rowIndex * unitY;

	let points = [
		[centerX,             centerY - hexSide],     // top
		[centerX + unitX / 2, centerY - hexSide / 2], // top right
		[centerX + unitX / 2, centerY + hexSide / 2], // bottom right
		[centerX,             centerY + hexSide],     // bottom
		[centerX - unitX / 2, centerY + hexSide / 2], // bottom left
		[centerX - unitX / 2, centerY - hexSide / 2], // top left
	];

	let ps = "";
	for (let p = 0; p < points.length; p++) {
		// Add margin to all coordinates to avoid negative coordinates
		ps += (points[p][0] + (unitX / 2 + margin)) + "," + (points[p][1] + (hexSide + margin)) + " ";
	}
	return ps;
}
