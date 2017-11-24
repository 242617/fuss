const LocalStorageKey = "242617/torture:address"

let address = () => document.getElementById("addressValue").value;

$(() => {

	document.getElementById("addressBtn").addEventListener("click", () => {
		window.localStorage.setItem(LocalStorageKey, document.getElementById("addressValue").value)
	});

	let enabled = document.getElementById("enabled");
	enabled.addEventListener("change", event => {
		fetch(address(), { method: "PUT", body: JSON.stringify({ enabled: event.target.checked }) })
		.then(console.log)
		.catch(console.error);
	});

	let slider  = $("#slider");
	slider.slider({
		animate: true,
		range: "min",
		min: 0,
		value: 100,
		step: 10,

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

			fetch(address(), { method: "PUT", body: JSON.stringify({ volume: value }) })
			.then(console.log)
			.catch(console.error);
		},
		// start: (event, ui) => {},
		// stop: (event, ui) => {},
	});

	fetch(address(), { method: "GET" })
	.then(response => response.json())
	.then(status => {
		enabled.checked = status.enabled;
		slider.slider("value", status.volume);
		new Switchery(document.getElementById("enabled"))
	})
	.catch(console.error);

	document.getElementById("addressValue").value = window.localStorage.getItem(LocalStorageKey)

});