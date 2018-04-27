function visualize_mcts() {
	return new Vue({
		el: '#mcts',
		data: {
			json: mcst_json.tree.root,
			showN: true,
			showQ: true,
			showGrid: true,
			showLastPlayer: true
		},
		props: {
			model: Object
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
	let app = visualize_mcts();
};
