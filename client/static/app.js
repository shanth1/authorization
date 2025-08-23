document.addEventListener("DOMContentLoaded", async () => {
	console.log("DOM loaded");

	UI.showLoader();

	try {
		const config = await Config.load();
		window.config = config;

		console.log("App initialized with config:", config);

		UI.createButtons();
		UI.showApp();
		UI.hideLoader();
	} catch (error) {
		console.error("Failed to initialize application:", error);
		UI.showError("Failed to load application. Please refresh the page.");
	}
});
