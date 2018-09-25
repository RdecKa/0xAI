function abSearchTree() {
	return new Vue({
		el: '#absearch',
		data: {
			size: 11,
			state: null,
			model: null
		},
		methods: {
			setJSON: function (dataJSON) {
				this.size = dataJSON.root.value.size;
				this.state = dataJSON.root.value.state;
				this.model = dataJSON.root.value.rootAB;
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
					<span>val: {{ model.value.val }}, </span>
					<span v-if="model.value.action">act:
						({{ model.value.action.x }},
						{{ model.value.action.y }},
						{{ model.value.action.c }})
					</span>
					<span>({{ model.value.comment }})</span>
				</span>
				<div v-if="showBoard" v-html="boardHTML" class="simple-board"></div>
			</div>

			<ul v-show="open" v-if="!isLeaf">
				<node
					v-for="(model, index) in model.children"
					:key="index"
					:size="size"
					:state="state"
					:model="model"
				></node>
			</ul>
		</li>`,
	props: {
		size: Number,
		state: Object,
		model: Object
	},
	data: function () {
		return {
			open: false,
			showBoard: false,
		}
	},
	computed: {
		isLeaf: function () {
			return !this.model.children || !this.model.children.length;
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
		getPreviousActions: function(obj, depth) {
			if (obj.model.value.action == null) {
				return [];
			}
			let prevActions = obj.$parent.$options.methods.getPreviousActions(obj.$parent, depth + 1);
			prevActions.push(obj.model.value.action);
			return prevActions;
		},
		drawBoard: function () {
			let prevActions = this.getPreviousActions(this, 0);

			let r = this.drawTopBottomRow();

			for (let row = 0; row < this.size; row++) {
				r += "<div>";

				r += this.indent(row);

				// Draw left blue column
				r += "<span class=\"cell blue\"></span>";

				let rowN = this.state.grid[row];
				for (let col = 0; col < this.size; col++) {
					let c = rowN & 3;
					switch (c) {
						case 0:
							let set = false
							for (let a = 0; a < prevActions.length; a++) {
								if (prevActions[a].x == col && prevActions[a].y == row) {
									switch (prevActions[a].c) {
										case "r":
											r += "<span class=\"cell red\"></span>";
											set = true;
											break;
										case "b":
											r += "<span class=\"cell blue\"></span>";
											set = true;
											break;
										default:
											console.log("Invalid color: " + prevActions[a].c);
									}
								}
							}
							if (!set) {
								r += "<span class=\"cell empty\"></span>";
							}
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
