

let address = () => document.getElementById("address").value;

$(() => {

	let enabled = document.getElementById("enabled");
	enabled.addEventListener("change", event => {
		fetch(address(), { method: "PUT", body: JSON.stringify({ enabled: event.target.checked }) })
		.then(console.log)
		.catch(console.error);
	});
	new Switchery(enabled);



	let slider  = $("#slider");

	slider.slider({
		animate: true,
		range: "min",
		min: 0,
		value: 100,
		step: 10,

		slide: (event, ui) => {

			let value  = slider.slider('value'),

			volume = $('.volume');
			if (value <= 5) { 
				volume.css('background-position', '0 0');
			} else if (value <= 25) {
				volume.css('background-position', '0 -25px');
			} else if (value <= 75) {
				volume.css('background-position', '0 -50px');
			} else {
				volume.css('background-position', '0 -75px');
			};

			fetch(address(), { method: "PUT", body: JSON.stringify({ volume: value }) })
			.then(console.log)
			.catch(console.error);

		},
		// start: (event, ui) => {},
		// stop: (event, ui) => {},

	});


});