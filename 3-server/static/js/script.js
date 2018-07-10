window.onload = function() {
	const socket = new WebSocket("ws://localhost:8080/ws/");

	socket.addEventListener("open", function (event) {
		socket.send("Hello Server!");
	});

	socket.addEventListener("message", function (event) {
		console.log("Message from server ", event.data);
	});
};
