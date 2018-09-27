let abTreeJSON;

function abSearchTree() {
	return new Vue({
		el: '#absearch',
		data: {
			initialised: false,
		},
		methods: {
			setJSON: function (dataJSON) {
				this.initialised = false;
				abTreeJSON = dataJSON.root;
				this.initialised = true;
			}
		}
	});
}

function getNodeFromIndices(treeJSON, indices) {
	if (treeJSON == undefined || treeJSON == null) {
		return null;
	}

	if (indices.length <= 0) {
		return treeJSON;
	}

	let subTree = treeJSON.children[indices[0]]
	return getNodeFromIndices(subTree, indices.slice(1));
}

Vue.component('node', {
	template: `<li v-if="initialised">
			<div class="node"
			:class="{extendable : !isLeaf}"
			@click="toggle"
			@mouseover="showBoard = true"
			@mouseleave="showBoard = false">
				<span class="info">
					<span v-if="!isLeaf" class="expand-sign">[{{ open ? '-' : '+' }}]</span>
					<span>value: {{ node.value.val }}, </span>
					<span>lastPlayer: {{ node.value.state.lastPlayer }}, </span>
					<!--<span>grid: {{ node.value.state.grid }} </span>-->
					<span>({{ node.value.comment }})</span>
				</span>
				<div v-if="showBoard" v-html="boardHTML" class="simple-board"></div>
			</div>

			<ul v-show="open" v-if="!isLeaf">
				<node
					v-for="(c, i) in children"
					:key="i"
					:modelindex="addIndexToModelindex(modelindex, i, c)"
					:initialised="initialised"
				></node>
			</ul>
		</li>`,
	props: {
		modelindex: Array,
		initialised: Boolean,
	},
	data: function () {
		return {
			open: false,
			showBoard: false,
		}
	},
	computed: {
		isLeaf: function () {
			if (!this.initialised) {
				return false;
			}
			let node = getNodeFromIndices(abTreeJSON, this.modelindex)
			return node == null || !node.children || !node.children.length;
		},
		size: function () {
			if (!this.initialised) {
				return 0;
			}
			return this.node.value.state.grid.length;
		},
		boardHTML: function () {
			if (!this.initialised) {
				return "";
			}
			return this.drawBoard();
		},
		node: function () {
			return getNodeFromIndices(abTreeJSON, this.modelindex);
		},
		children: function () {
			if (!this.initialised) {
				return [];
			}
			if (!this.open) {
				// For the sake of RAM well-being
				return [];
			}
			return this.node.children;
		},
	},
	methods: {
		toggle: function () {
			if (!this.isLeaf) {
				this.open = !this.open;
			}
		},
		drawBoard: function () {
			let r = this.drawTopBottomRow();

			for (let row = 0; row < this.size; row++) {
				r += "<div>";

				r += this.indent(row);

				// Draw left blue column
				r += drawCellWithClass("cell blue");

				let rowN = this.node.value.state.grid[row];
				for (let col = 0; col < this.size; col++) {
					let c = rowN & 3;
					switch (c) {
						case 0:
							r += drawCellWithClass("cell empty");
							break;
						case 1:
							r += drawCellWithClass("cell red");
							break;
						case 2:
							r += drawCellWithClass("cell blue");
							break;
						default:
							r += drawCellWithClass("cell undef");
					}
					rowN = rowN >> 2;
				}

				// Draw right blue column
				r += drawCellWithClass("cell blue");

				r += "</div>";
			}

			r += this.indent(this.size);
			r += this.drawTopBottomRow();

			return r;
		},
		drawTopBottomRow() {
			// Draw top/bottom red row with two corner cells
			let r = "";
			r += drawCellWithClass("cell violet");
			for (let col = 0; col < this.size; col++) {
				r += drawCellWithClass("cell red");
			}
			r += drawCellWithClass("cell violet");
			return r;
		},
		indent(indent) {
			let r = "";
			for (let i = 0; i <= indent; i++) {
				r += drawCellWithClass("indent");
			}
			return r;
		},
		addIndexToModelindex: function (indices, indexToAdd) {
			return indices.concat([indexToAdd]);
		}
	}
})

function drawCellWithClass(cl) {
	return "<span class=\"" + cl + "\"></span>";
}
