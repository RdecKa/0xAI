const colors = Object.freeze({NONE: ".", RED: "r", BLUE: "b"});

window.onload = function() {
	const socket = new WebSocket("ws://localhost:8080/ws/");
	let hexgrid = hexGrid(socket);
};

/**
 * Returns an action to be performed.
 * @param grid    current game board
 * @param myColor color of the player
 */
function nextMove(grid, myColor) {
	// TODO: random for now
	for (let i = 0; i < grid.length; i++) {
		for (let j = 0; j < grid[0].length; j++) {
			if (grid[i][j] == colors.NONE) {
				grid[i].splice(j, 1, myColor)
				return encodeMove({x: j, y: i});
			}
		}
	}
}

/**
 * Prints the given 2D grid.
 * @param grid 2D grid to be printed.
 */
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

/**
 * Returns the string representation of a color, used for setting class of
 * cells.
 * @param color color to be represented as a string
 */
function getCellColorName(color) {
	switch (color) {
		case colors.RED:
			return "red";
		case colors.BLUE:
			return "blue";
		default:
			return "empty";
	}
}

/**
 * Formats the given objects in a way that can be sent to server.
 * @param moveObj move object to be encoded and sent
 */
function encodeMove(moveObj) {
	return moveObj.x.toString() + "," + moveObj.y.toString();
}

/**
 * Reads coordinates and color from the string representing a move (received
 * from server). Example of a string received: 'r: (2, 3)'.
 * @param moveString
 */
function decodeMove(moveString) {
	let color;
	switch (moveString.charAt(0)) {
		case 'r':
			color = colors.RED;
			break;
		case 'b':
			color = colors.BLUE;
			break;
		default:
			console.log("INVALID COLOR '" + moveString.charAt(0) + "'");
			color = colors.NONE;
	}
	let coords = moveString.substring(4, moveString.length - 1).split(", ");
	return {x: parseInt(coords[0]), y: parseInt(coords[1]), c: color};
}
