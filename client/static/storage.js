const Storage = (function () {
	function saveConfig(config) {
		localStorage.setItem("authUrl", config.authUrl);
		localStorage.setItem("clientId", config.clientId);
		console.log("Configuration saved to localStorage");
	}

	function getConfig() {
		return {
			authUrl: localStorage.getItem("authUrl"),
			clientId: localStorage.getItem("clientId"),
		};
	}

	function hasConfig() {
		return (
			localStorage.getItem("authUrl") && localStorage.getItem("clientId")
		);
	}

	function clearConfig() {
		localStorage.removeItem("authUrl");
		localStorage.removeItem("clientId");
		console.log("Configuration cleared from localStorage");
	}

	return {
		saveConfig,
		getConfig,
		hasConfig,
		clearConfig,
	};
})();
