// English comment: Script for AppFrontend: Handle login initiation, callback, token management, user info, refresh, logout.

document.addEventListener("DOMContentLoaded", () => {
	const loginBtn = document.getElementById("login-btn");
	const logoutBtn = document.getElementById("logout-btn");
	const userInfoDiv = document.getElementById("user-info");
	const userDataPre = document.getElementById("user-data");

	// Hardcoded for demo: In production, from config
	const clientId = "example_client_id";
	const redirectUri = "http://localhost:8081/callback"; // Must match registered in backend
	const authUrl = "http://localhost:8082"; // AuthFrontend URL

	// Check if callback (after login)
	const params = new URLSearchParams(window.location.search);
	if (params.has("code") && params.has("state")) {
		handleCallback(params.get("code"), params.get("state"));
	} else {
		checkIfLoggedIn();
	}

	// Login button: Generate params and redirect
	loginBtn.addEventListener("click", async () => {
		const state = generateRandomString(32);
		const nonce = generateRandomString(32);
		const codeVerifier = generateRandomString(43); // 43-128 chars for PKCE
		const hash = await sha256(codeVerifier);
		const codeChallenge = base64UrlEncode(hash);
		debugger;

		// Save state and codeVerifier to localStorage for callback check
		localStorage.setItem("state", state);
		localStorage.setItem("code_verifier", codeVerifier);

		const authorizeUrl = `${authUrl}?response_type=code&client_id=${encodeURIComponent(clientId)}&redirect_uri=${encodeURIComponent(redirectUri)}&scope=openid profile email&state=${encodeURIComponent(state)}&code_challenge=${encodeURIComponent(codeChallenge)}&code_challenge_method=S256&nonce=${encodeURIComponent(nonce)}`;
		console.log("URL:", authorizeUrl);
		window.location.href = authorizeUrl;
	});

	// Logout button
	logoutBtn.addEventListener("click", () => {
		fetch("/logout", { method: "POST" })
			.then(() => {
				localStorage.clear();
				userInfoDiv.style.display = "none";
				loginBtn.style.display = "block";
				window.location.href = "/"; // Redirect to clean URL
			})
			.catch((error) => console.error("Logout error:", error));
	});

	// Function to handle callback: Exchange code for tokens via AppBackend
	function handleCallback(code, receivedState) {
		const savedState = localStorage.getItem("state");
		const codeVerifier = localStorage.getItem("code_verifier");

		if (receivedState !== savedState) {
			alert("Invalid state: CSRF attack possible");
			return;
		}

		const formData = new FormData();
		formData.append("code", code);
		formData.append("code_verifier", codeVerifier);

		fetch("/token", {
			method: "POST",
			body: formData,
		})
			.then((response) => response.json())
			.then((data) => {
				// Assume backend sets secure cookies; for demo, save to localStorage
				localStorage.setItem("access_token", data.access_token);
				localStorage.setItem("id_token", data.id_token);
				localStorage.setItem("refresh_token", data.refresh_token);

				// Validate nonce in id_token (parse JWT)
				const idTokenClaims = parseJwt(data.id_token);
				// Check nonce matches, but nonce was sent, assume backend validated

				localStorage.removeItem("state");
				localStorage.removeItem("code_verifier");

				displayUserInfo();
			})
			.catch((error) => console.error("Token exchange error:", error));
	}

	// Check if logged in (has tokens)
	function checkIfLoggedIn() {
		const accessToken = localStorage.getItem("access_token");
		if (accessToken) {
			displayUserInfo();
		}
	}

	// Display user info: Fetch /userinfo with access_token
	function displayUserInfo() {
		let accessToken = localStorage.getItem("access_token");
		fetch("/userinfo", {
			headers: { Authorization: `Bearer ${accessToken}` },
		})
			.then((response) => {
				if (response.status === 401) {
					// Expired, refresh
					return refreshTokens().then((newToken) => {
						accessToken = newToken;
						return fetch("/userinfo", {
							headers: { Authorization: `Bearer ${accessToken}` },
						});
					});
				}
				return response;
			})
			.then((response) => response.json())
			.then((data) => {
				userDataPre.textContent = JSON.stringify(data, null, 2);
				userInfoDiv.style.display = "block";
				loginBtn.style.display = "none";
			})
			.catch((error) => console.error("User info error:", error));
	}

	// Refresh tokens via /refresh
	function refreshTokens() {
		const refreshToken = localStorage.getItem("refresh_token");
		if (!refreshToken) {
			throw new Error("No refresh token");
		}

		const formData = new FormData();
		formData.append("refresh_token", refreshToken);

		return fetch("/refresh", {
			method: "POST",
			body: formData,
		})
			.then((response) => response.json())
			.then((data) => {
				localStorage.setItem("access_token", data.access_token);
				localStorage.setItem("refresh_token", data.refresh_token); // Rotated
				return data.access_token;
			});
	}

	// Helpers
	function generateRandomString(length) {
		const array = new Uint8Array(length);
		window.crypto.getRandomValues(array);
		return Array.from(array, (byte) =>
			byte.toString(16).padStart(2, "0"),
		).join("");
	}

	function sha256(str) {
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

	function parseJwt(token) {
		const base64Url = token.split(".")[1];
		const base64 = base64Url.replace(/-/g, "+").replace(/_/g, "/");
		const jsonPayload = decodeURIComponent(
			atob(base64)
				.split("")
				.map(
					(c) =>
						"%" + ("00" + c.charCodeAt(0).toString(16)).slice(-2),
				)
				.join(""),
		);
		return JSON.parse(jsonPayload);
	}
});
