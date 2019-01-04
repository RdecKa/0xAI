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
			boardSize: 11,
			numGames: 4
		},
		methods: {
			clickPlayHandler: function () {
				let newLocation = "http://localhost:8080/play/"
					+ "?red=" + this.selection.Red.type
					+ "&blue=" + this.selection.Blue.type
					+ "&watch=" + this.watchInBrowser
					+ "&size=" + this.boardSize
					+ "&numgames=" + this.numGames;

				if (this.selection.Red.type != "human" && this.selection.Red.type != "rand") {
					newLocation += "&redtime=" + this.selection.Red.time
				}
				if (this.selection.Blue.type != "human" && this.selection.Blue.type != "rand") {
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
	template: `<div :class="'col col-2 player-' + color">
			<div class="container">
				<h2>{{ color }} player</h2>
				<input type="radio" :id="'human-' + color" :name="color" value="human" v-model="player" @change="selectionChange" />
				<label :for="'human-' + color">Human</label>
				<br>
				<input type="radio" :id="'rand-' + color" :name="color" value="rand" v-model="player" @change="selectionChange" />
				<label :for="'rand-' + color">RAND</label>
				<br>
				<input type="radio" :id="'mcts-'  + color" :name="color" value="mcts"  v-model="player" @change="selectionChange" />
				<label :for="'mcts-'  + color">MCTS</label>
				<input type="number" min="1" :id="'time-mcts-' + color" v-model="time.mcts" @change="selectionChange">
				<label :for="'time-mcts-' + color">seconds</label>
				<br>
				<input type="radio" :id="'abDT-'  + color" :name="color" value="abDT"  v-model="player" @change="selectionChange" />
				<label :for="'abDT-'  + color">ABDL</label>
				<input type="number" min="1" :id="'time-abDT-' + color" v-model="time.abDT" @change="selectionChange">
				<label :for="'time-abDT-' + color">seconds</label>
				<br>
				<input type="radio" :id="'abLR-'  + color" :name="color" value="abLR"  v-model="player" @change="selectionChange" />
				<label :for="'abLR-'  + color">ABLR</label>
				<input type="number" min="1" :id="'time-abLR-' + color" v-model="time.abLR" @change="selectionChange">
				<label :for="'time-abLR-' + color">seconds</label>
				<br>
				<input type="radio" :id="'hybrid-'  + color" :name="color" value="hybrid"  v-model="player" @change="selectionChange" />
				<label :for="'hybrid-'  + color">HYBR</label>
				<input type="number" min="1" :id="'time-hybrid-' + color" v-model="time.hybrid" @change="selectionChange">
				<label :for="'time-hybrid-' + color">seconds</label>
			</div>
		</div>`,
	props: ["color"],
	data: function () {
		return {
			player: null,
			time: {mcts: 1, abDT: 1, abLR: 1, hybrid: 1},
		}
	},
	methods: {
		selectionChange: function () {
			this.$emit("selection-change", {color: this.color, type: this.player, time: this.time[this.player]});
		}
	}
})
