Vue.component("hex-cell", {
	template: `<polygon	:points='points'
		class='cell'
		:class='[cellColor, cellActive]'
		@mouseover='mouseChange()'
		@mouseleave='mouseChange()'
		></polygon>`,
	props: ["x", "y", "color", "playercolor"],
	data: function () {
		return {
			active: false
		}
	},
	computed: {
		points: function () {
			return getPointsForPolygon(this.y, this.x);
		},
		cellColor: function () {
			if (this.color == colors.NONE && this.active) {
				// Cell is hovered, color it with player's color
				return getCellColorName(this.playercolor);
			}
			return getCellColorName(this.color);
		},
		cellActive: function () {
			if (this.active && this.color == colors.NONE && this.playercolor != colors.NONE) {
				return "active";
			};
		}
	},
	methods: {
		mouseChange: function () {
			this.active = !this.active;
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
