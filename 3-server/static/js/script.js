const colors = Object.freeze({NONE: ".", RED: "r", BLUE: "b"});

window.onload = function() {
	const socket = new WebSocket("ws://localhost:8080/ws/");

	let myColor = colors.NONE;
	let grid;

	socket.addEventListener("open", function (event) {
		console.log("Connection established.");
	});

	socket.addEventListener("message", function (event) {
		let m = event.data;
		console.log("Message from server:", m);


		let ms = m.split(" ");
		switch (ms[0]) {
			case "INIT":
				console.log("Init");
				let s = ms[1].split(":");
				grid = creatHexeGrid(parseInt(s[1]));

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
					socket.send("READY");
				} else {
					socket.send("ERROR");
				}
				break;
			case "MOVE":
				console.log("Move");
				if (ms[1] !== "<nil>") {
					receiveMove(grid, m.substring(5))
				}
				socket.send(nextMove(grid, myColor));
				printGrid(grid);
				break;
			case "END":
				console.log("End");
				receiveMove(grid, m.substring(6));
				printGrid(grid);
				myColor = colors.NONE;
				break;
			default:
				console.log("Unknown message")
		}
	});
};

function creatHexeGrid(size) {
	let grid = [];
	for (let i = 0; i < size; i++) {
		grid[i] = [];
		for (let j = 0; j < size; j++) {
			grid[i][j] = colors.NONE;
		}
	}
	return grid;
};

function receiveMove(grid, move) {
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

	grid[parseInt(coords[1])][parseInt(coords[0])] = color;
}

function nextMove(grid, myColor) {
	for (let i = 0; i < grid.length; i++) {
		for (let j = 0; j < grid[0].length; j++) {
			if (grid[i][j] == colors.NONE) {
				grid[i][j] = myColor;
				return j.toString() + "," + i.toString();
			}
		}
	}
};

function opponent(myColor) {
	if (myColor == colors.BLUE) return colors.RED;
	if (myColor == colors.RED) return colors.BLUE;
	return colors.NONE
}

function printGrid(grid) {
	for (let i = 0; i < grid.length; i++) {
		let line = "";
		for (let s = 0; s < i; s++) {
			line += " ";
		}
		for (let j = 0; j < grid[i].length; j++) {
			switch (grid[i][j]) {
				case colors.RED:
					line += "r ";
					break;
				case colors.BLUE:
					line += "b ";
					break;
				default:
				line += ". ";
			}
		}
		console.log(line);
	}
}
