const Auth = (function () {
	function handleLogin() {
		const config = Storage.getConfig();
		console.log("Login with:", config);
	}

	function handleLogout() {
		console.log("Logout");
	}

	return {
		handleLogin,
		handleLogout,
	};
})();
