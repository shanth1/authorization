const Auth = (function () {
	async function handleLogin() {
		const config = window.config;

		const clientId = config.clientId;
		const authUrl = config.authUrl;

		const redirectUri = window.location.href;
		const state = generateRandomString(32);
		const nonce = generateRandomString(32);
		const codeVerifier = generateRandomString(43);
		const hash = await sha256(codeVerifier);
		const codeChallenge = base64UrlEncode(hash);

		localStorage.setItem("state", state);
		localStorage.setItem("code_verifier", codeVerifier);

		const authorizeUrl = `${authUrl}/autorize?response_type=code&client_id=${encodeURIComponent(clientId)}&redirect_uri=${encodeURIComponent(redirectUri)}&scope=openid profile email&state=${encodeURIComponent(state)}&code_challenge=${encodeURIComponent(codeChallenge)}&code_challenge_method=S256&nonce=${encodeURIComponent(nonce)}`;
		window.location.href = authorizeUrl;
	}

	function handleLogout() {
		localStorage.clear();
	}

	return {
		handleLogin,
		handleLogout,
	};
})();

function generateRandomString(length) {
	const array = new Uint8Array(length);
	window.crypto.getRandomValues(array);
	return Array.from(array, (byte) => byte.toString(16).padStart(2, "0")).join(
		"",
	);
}

async function sha256(str) {
	const buffer = new TextEncoder().encode(str);
	return crypto.subtle
		.digest("SHA-256", buffer)
		.then((hash) => new Uint8Array(hash));
}

function base64UrlEncode(array) {
	return btoa(String.fromCharCode.apply(null, array))
		.replace(/\+/g, "-")
		.replace(/\//g, "_")
		.replace(/=+$/, "");
}
