function visualize_mcts() {
	return new Vue({
		el: '#mcts',
		data: {
			json: mcs_json.tree.root,
		},
		props: {
			model: Object
		}
	});
}

Vue.component('item', {
	template: '#item-template',
	props: {
		model: Object
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
