async function loadConfig() {
	try {
		const response = await fetch("/api/config");
		const config = await response.json();

		return config;
	} catch (error) {
		console.error("Failed to load configuration:", error);

		// default config
		return {
			authUrl: "http://localhost:8000",
			clientID: "",
		};
	}
}

loadConfig().then((config) => {
	initApp(config);
});

function initApp(config) {
	console.log("App initialized with config:", config);
}
