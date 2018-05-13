function visualize_mcts() {
	return new Vue({
		el: '#mcts',
		data: {
			json: mcst_json.tree.root,
			showN: true,
			showQ: true,
			showGrid: true,
			showLastPlayer: true,
			sortBy: "n"
		},
		props: {
			model: Object
		},
		watch: {
			sortBy: function(val) {
				switch(val) {
					case "q":
						sortJSON(this.json, compareByQ);
						break;
					case "n":
						sortJSON(this.json, compareByN);
						break;
				}
			}
		}
	});
}

Vue.component('item', {
	template: '#item-template',
	props: {
		model: Object,
		showN: Boolean,
		showQ: Boolean,
		showGrid: Boolean,
		showLastPlayer: Boolean,
	},
	data: function () {
		return {
			open: false,
			showBoard: false,
			boardHTML: null
		}
	},
	created() {
		this.boardHTML = this.drawBoard();
	},
	computed: {
		isGoal: function () {
			return !this.model.children || !this.model.children.length
		},
		size: function () {
			return this.model.value.state.grid.length;
		}
	},
	methods: {
		toggle: function () {
			if (!this.isGoal) {
				this.open = !this.open
			}
		},
		drawBoard: function () {
			let r = this.drawTopBottomRow();

			for (let row = 0; row < this.size; row++) {
				r += "<div>"

				r += this.indent(row);

				// Draw left blue column
				r += "<span class=\"cell blue\"></span>";

				let rowN = this.model.value.state.grid[row];
				for (let col = 0; col < this.size; col++) {
					let c = rowN & 3;
					switch (c) {
						case 0:
							r += "<span class=\"cell empty\"></span>";
							break;
						case 1:
							r += "<span class=\"cell red\"></span>";
							break;
						case 2:
							r += "<span class=\"cell blue\"></span>";
							break;
						default:
							r += "<span class=\"cell undef\"></span>";
					}
					rowN = rowN >> 2;
				}

				// Draw right blue column
				r += "<span class=\"cell blue\"></span>";

				r += "</div>";
			}

			r += this.indent(this.size);
			r += this.drawTopBottomRow();

			return r;
		},
		drawTopBottomRow() {
			// Draw top/bottom red row with two corner cells
			let r = "";
			r += "<span class=\"cell violet\"></span>";
			for (let col = 0; col < this.size; col++) {
				r += "<span class=\"cell red\"></span>";
			}
			r += "<span class=\"cell violet\"></span>";
			return r;
		},
		indent(indent) {
			let r = "";
			for (let i = 0; i <= indent; i++) {
				r += "<span class=\"indent\"></span>";
			}
			return r;
		}
	}
});

window.onload = function() {
	sortJSON(mcst_json.tree.root, compareByN);
	let app = visualize_mcts();
};

function sortJSON(json, sortBy) {
	sortChildren(json.children, sortBy);
	for (let c in json.children) {
		sortJSON(json.children[c], sortBy);
	}
}

function sortChildren(children, sortBy) {
	children.sort(sortBy);
}

function compareByN(a, b) {
	return b.value.N - a.value.N;
}

function compareByQ(a, b) {
	return b.value.Q - a.value.Q;
}
