const LocalStorageKey = "242617/torture:address"

let getAddress = () => document.getElementById("addressValue").value;
let setAddress = address => document.getElementById("addressValue").value = address;

$(() => {
	setAddress(window.localStorage.getItem(LocalStorageKey));

	let left = document.getElementById("leftValue");
	let right = document.getElementById("rightValue");

	document.getElementById("enabled").addEventListener("change", event => {
		fetch(getAddress(), { method: "PUT", body: JSON.stringify({ enabled: event.target.checked }) })
		.then(console.log)
		.catch(console.error);
	});

	document.getElementById("frequencyBtn").addEventListener("click", event => {
		console.log("here");
		fetch(getAddress(), { method: "PUT", body: JSON.stringify({
			left: parseInt(left.value) || 0,
			right: parseInt(right.value) || 0
		}) })
		.then(console.log)
		.catch(console.error);
	});

	document.getElementById("addressBtn").addEventListener("click", () => {
		window.localStorage.setItem(LocalStorageKey, document.getElementById("addressValue").value);
	});

	let slider  = $("#slider");
	slider.slider({
		animate: true,
		range: "min",
		min: 0,
		value: 100,
		slide: (event, ui) => {
			let value  = slider.slider("value"),

			volume = $(".volume");
			if (value <= 5) { 
				volume.css("background-position", "0 0");
			} else if (value <= 25) {
				volume.css("background-position", "0 -25px");
			} else if (value <= 75) {
				volume.css("background-position", "0 -50px");
			} else {
				volume.css("background-position", "0 -75px");
			};

			fetch(getAddress(), { method: "PUT", body: JSON.stringify({ volume: value }) })
			.then(console.log)
			.catch(console.error);
		},
		// start: (event, ui) => {},
		// stop: (event, ui) => {},
	});

	fetch(getAddress(), { method: "GET" })
	.then(response => response.json())
	.then(status => {
		enabled.checked = status.enabled;
		slider.slider("value", status.volume);
		left.value = status.left;
		right.value = status.right;
		new Switchery(document.getElementById("enabled"))
	})
	.catch(console.error);

});