const Config = (function () {
	async function load() {
		try {
			const response = await fetch("/api/config");
			if (!response.ok) {
				throw new Error(`HTTP error! status: ${response.status}`);
			}
			const config = await response.json();
			console.log("Configuration loaded from server");
			return config;
		} catch (error) {
			console.error("Failed to load configuration:", error);

			return {
				authUrl: "http://localhost:8000",
				clientId: "",
			};
		}
	}

	return {
		load,
	};
})();
