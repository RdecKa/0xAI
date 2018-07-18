Vue.component("hex-cell", {
	props: ["x", "y", "color"],
	template: "<polygon	:points='points' :class='cellClass'></polygon>",
	computed: {
		points: function () {
			return getPointsForPolygon(this.y, this.x);
		},
		cellClass: function () {
			let c = "cell ";
			switch (this.color) {
				case colors.RED:
					return c + "red";
				case colors.BLUE:
					return c + "blue";
				default:
					return c + "empty";
			}
		}
	}
})

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
