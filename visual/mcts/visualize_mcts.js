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
					case "n":
						sortJSON(this.json, compareByN);
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
			open: false
		}
	},
	computed: {
		isGoal: function () {
			return !this.model.children || !this.model.children.length
		}
	},
	methods: {
		toggle: function () {
			if (!this.isGoal) {
				this.open = !this.open
			}
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
