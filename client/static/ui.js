const UI = (function () {
	function showLoader() {
		document.getElementById("loader").classList.remove("hidden");
	}

	function hideLoader() {
		document.getElementById("loader").classList.add("hidden");
	}

	function showApp() {
		document.getElementById("app-container").classList.remove("hidden");
	}

	function createButtons() {
		const appContainer = document.getElementById("app-container");

		const loginBtn = document.createElement("button");
		loginBtn.id = "login-btn";
		loginBtn.textContent = "Login";
		loginBtn.addEventListener("click", Auth.handleLogin);

		const userInfo = document.createElement("div");
		userInfo.id = "user-info";
		userInfo.innerHTML = `
            <h2>User info</h2>
            <pre id="user-data"></pre>
            <button id="logout-btn">Logout</button>
        `;

		appContainer.appendChild(loginBtn);
		appContainer.appendChild(userInfo);

		document
			.getElementById("logout-btn")
			.addEventListener("click", Auth.handleLogout);

		console.log("UI buttons created");
	}

	return {
		showLoader,
		hideLoader,
		showApp,
		createButtons,
	};
})();
