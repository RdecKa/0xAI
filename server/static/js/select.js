window.onload = function() {
	let sp = selectPlayers();
};

function selectPlayers() {
	return new Vue({
		el: '#select-player-wrap',
		data: {
			buttonActive: false,
			message: "Select both players, please.",
			selection: {Red: {type:null, time:0}, Blue: {type:null, time:0}},
			watchInBrowser: true,
			watchInBrowserDisabled: false,
			boardSize: 7,
			numGames: 10
		},
		methods: {
			clickPlayHandler: function () {
				let newLocation = "http://localhost:8080/play/"
					+ "?red=" + this.selection.Red.type
					+ "&blue=" + this.selection.Blue.type
					+ "&watch=" + this.watchInBrowser
					+ "&size=" + this.boardSize
					+ "&numgames=" + this.numGames;

				if (this.selection.Red.type != "human") {
					newLocation += "&redtime=" + this.selection.Red.time
				}
				if (this.selection.Blue.type != "human") {
					newLocation += "&bluetime=" + this.selection.Blue.time
				}

				window.location.href = newLocation;
			},
			onSelectionChange: function (event) {
				this.selection[event.color].type = event.type;
				this.selection[event.color].time = event.time;
				if (this.selection.Red  == null || this.selection.Red.type  == null ||
					this.selection.Blue == null || this.selection.Blue.type == null) {
					this.buttonActive = false;
					this.message = "Select both players, please.";
					return;
				}
				if (this.selection.Red.type == "human" && this.selection.Blue.type == "human") {
					this.buttonActive = false;
					this.message = "Sorry, this combination is not supported.";
					return;
				}

				this.buttonActive = true;
				this.message = "Let's play!";

				if (this.selection.Red.type == "human" || this.selection.Blue.type == "human") {
					this.watchInBrowser = true;
					this.watchInBrowserDisabled = true;
				} else {
					this.watchInBrowserDisabled = false;
				}
			}
		}
	});
}

Vue.component("select-player", {
	template: `<div class="col col-2">
			<h2>{{ color }} player</h2>
			<input type="radio" :id="'human-' + color" :name="color" value="human" v-model="player" @change="selectionChange" />
			<label :for="'human-' + color">Human</label>
			<br>
			<input type="radio" :id="'mcts-'  + color" :name="color" value="mcts"  v-model="player" @change="selectionChange" />
			<label :for="'mcts-'  + color">Computer (MCTS)</label>
			<input type="number" min="1" :id="'time-mcts-' + color" v-model="time.mcts" @change="selectionChange">
			<label :for="'time-mcts-' + color">seconds</label>
			<br>
			<input type="radio" :id="'ab-'  + color" :name="color" value="ab"      v-model="player" @change="selectionChange" />
			<label :for="'ab-'  + color">Computer (AB)</label>
			<input type="number" min="1" :id="'time-ab-' + color" v-model="time.ab" @change="selectionChange">
			<label :for="'time-ab-' + color">seconds</label>
		</div>`,
	props: ["color"],
	data: function () {
		return {
			player: null,
			time: {mcts: 1, ab: 1}
		}
	},
	methods: {
		selectionChange: function () {
			this.$emit("selection-change", {color: this.color, type: this.player, time: this.time[this.player]});
		}
	}
})
