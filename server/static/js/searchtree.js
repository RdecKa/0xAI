function abSearchTree() {
	return new Vue({
		el: '#absearch',
		data: {
			model: null,
		},
		methods: {
			setJSON: function (dataJSON) {
				this.model = dataJSON;
			}
		}
	});
}

Vue.component('node', {
	template: `<li>
			<div class="node"
			:class="{extendable : !isLeaf}"
			@click="toggle"
			@mouseover="showBoard = true"
			@mouseleave="showBoard = false">
				<span class="info">
					<span v-if="!isLeaf" class="expand-sign">[{{ open ? '-' : '+' }}]</span>
					<span>value: {{ model.value.val }}, </span>
					<span>lastPlayer: {{ model.value.state.lastPlayer }} </span>
					<!--<span>grid: {{ model.value.state.grid }} </span>-->
				</span>
				<div v-if="showBoard" v-html="boardHTML" class="simple-board"></div>
			</div>

			<ul v-show="open" v-if="!isLeaf">
				<node
				v-for="(model, index) in model.children"
				:key="index"
				:model="model"
				></node>
			</ul>
		</li>`,
	props: {
		model: Object,
	},
	data: function () {
		return {
			open: true,
			showBoard: false,
		}
	},
	computed: {
		isLeaf: function () {
			return !this.model.children || !this.model.children.length;
		},
		size: function () {
			return this.model.value.state.grid.length;
		},
		boardHTML: function () {
			return this.drawBoard();
		}
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
})
