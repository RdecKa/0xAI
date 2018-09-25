const hexSide = 30;
const unitX = 2 * Math.cos(30 * Math.PI / 180) * hexSide;
const unitY = (1 + Math.sin(30 * Math.PI / 180)) * hexSide;
const margin = hexSide * 0.1;

function hexGrid(socket, abSearchTree) {
	return new Vue({
		el: '#hex-grid',
		data: {
			socket: socket,
			size: 1,
			grid: [],
			myColor: colors.NONE,
			playersTurn: false,
			passivePlayer: false,
			abSearchTree: abSearchTree,
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
			onClickReceived: function (click) {
				sendMove(this, {x: click.x, y: click.y, c: this.myColor});
			},
			setIsMyTurn: function (isMyTurn) {
				this.playersTurn = isMyTurn;
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

	socket.addEventListener("close", function (event) {
		console.log("Connection closed.");
	});

	socket.addEventListener("message", function (event) {
		let m = event.data;
		//console.log("Message from server:", m);
		let response = respondToMessage(obj, m);
		if (response !== undefined) {
			socket.send(response);
		}
	});
}

/**
 * Reads the message and returns a response when needed.
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
			obj.setIsMyTurn(false);

			let c = ms[2].split(":")
			switch (c[1]) {
				case "r":
					obj.myColor = colors.RED;
					break;
				case "b":
					obj.myColor = colors.BLUE;
					break;
				default:
					obj.myColor = colors.NONE;
			}
			if (obj.myColor == colors.RED || obj.myColor == colors.BLUE) {
				obj.passivePlayer = false;
				return "READY";
			} else if (obj.myColor == colors.NONE) {
				obj.passivePlayer = true;
				return "READY PASSIVE";
			} else {
				return "ERROR";
			}
		case "MOVE":
			console.log("Move");
			if (ms[1] !== "<nil>") {
				receiveMove(obj, decodeMove(msg.substring(5)));
			}
			printGrid(obj.grid);
			if (!obj.passivePlayer) {
				obj.setIsMyTurn(true);
			}
			return;
		case "END":
			console.log("End");
			receiveMove(obj, decodeMove(msg.substring(6)));
			printGrid(obj.grid);
			obj.myColor = colors.NONE;
			setTimeout(function() { obj.socket.send("DONE"); }, 2000);
			return;
		case "ABJSON":
			console.log("Got JSON");
			let dataJSON = JSON.parse(msg.substring(7));
			obj.abSearchTree.setJSON(dataJSON);
			return;
		default:
			console.log("Unknown message: '" + msg + "'");
			return;
	}
}

/**
 * Decodes a move, changes the grid accordingly.
 * @param obj     Vue instance
 * @param moveObj object representing a move
 */
function receiveMove(obj, moveObj) {
	obj.grid[moveObj.y].splice(moveObj.x, 1, moveObj.c);
}

/**
 * Updates obj's state and sends selected move via obj's socket.
 * @param obj     Vue instance
 * @param moveObj move to be sent to server
 */
function sendMove(obj, moveObj) {
	receiveMove(obj, moveObj);
	obj.setIsMyTurn(false);
	obj.socket.send(encodeMove(moveObj));
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
