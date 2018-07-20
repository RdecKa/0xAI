window.onload = function() {
	let sp = selectPlayers();
};

function selectPlayers() {
	return new Vue({
		el: '#select-player-wrap',
		data: {
			buttonActive: false,
			message: "Select both players, please.",
			selection: {"Red": null, "Blue": null}
		},
		methods: {
			clickPlayHandler: function () {
				window.location.href = "http://localhost:8080/play/?red=" + this.selection["Red"] + "&blue=" + this.selection["Blue"];
			},
			onSelectionChange: function (event) {
				this.selection[event.color] = event.selection;
				if (this.selection["Red"] == null || this.selection["Blue"] == null) {
					this.buttonActive = false;
					this.message = "Select both players, please.";
				} else if (this.selection["Red"] == this.selection["Blue"]) {
					this.buttonActive = false;
					this.message = "Sorry, this combination is not supported.";
				} else {
					this.buttonActive = true;
					this.message = "Let's play!";
				}
			}
		}
	});
}

Vue.component("select-player", {
	template: `<div class="col-2">
			<h2>{{ color }} player</h2>
			<form>
				<input type="radio" :id="'human-' + color" :name="color" value="human" v-model="player" @change="selectionChange" />
				<label :for="'human-' + color">Human</label>
				<br>
				<input type="radio" :id="'mcts-'  + color" :name="color" value="mcts"  v-model="player" @change="selectionChange" />
				<label :for="'mcts-'  + color">Computer (MCTS)</label>
			</form>
		</div>`,
	props: ["color"],
	data: function () {
		return {
			player: null
		}
	},
	methods: {
		selectionChange: function () {
			this.$emit("selection-change", {color: this.color, selection: this.player});
		}
	}
})
