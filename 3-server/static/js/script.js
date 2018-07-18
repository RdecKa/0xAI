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
				return j.toString() + "," + i.toString();
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
