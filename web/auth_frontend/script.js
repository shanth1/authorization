// English comment: Script for AuthFrontend: Parse query params, fetch providers, render buttons and form, handle submissions.

document.addEventListener("DOMContentLoaded", () => {
	// Parse query parameters from URL
	const params = new URLSearchParams(window.location.search);
	const clientId = params.get("client_id");
	const redirectUri = params.get("redirect_uri");
	const state = params.get("state");
	const codeChallenge = params.get("code_challenge");
	const nonce = params.get("nonce");

	// Set hidden fields
	document.getElementById("client_id").value = clientId;
	document.getElementById("redirect_uri").value = redirectUri;
	document.getElementById("state").value = state;
	document.getElementById("code_challenge").value = codeChallenge;
	document.getElementById("nonce").value = nonce;

	// Fetch available providers from AuthBackend
	fetch("/providers")
		.then((response) => response.json())
		.then((providers) => {
			const providersDiv = document.getElementById("providers");
			providers.forEach((provider) => {
				if (provider === "login_password") {
					// Show login/password form
					document.getElementById("login-form").style.display =
						"block";
				} else {
					// Create button for external providers
					const btn = document.createElement("button");
					btn.className = "provider-btn";
					btn.textContent = `Войти через ${provider.charAt(0).toUpperCase() + provider.slice(1)}`;
					btn.onclick = () => initiateExternalLogin(provider);
					providersDiv.appendChild(btn);
				}
			});
		})
		.catch((error) => console.error("Error fetching providers:", error));

	// Handle login/password form submission
	document.getElementById("login-form").addEventListener("submit", (e) => {
		e.preventDefault();
		const formData = new FormData();
		formData.append("provider", "login_password");
		formData.append("login", document.getElementById("login").value);
		formData.append("password", document.getElementById("password").value);
		formData.append("client_id", clientId);
		formData.append("redirect_uri", redirectUri);
		formData.append("state", state);
		formData.append("code_challenge", codeChallenge);
		formData.append("nonce", nonce);

		fetch("/authorize", {
			method: "POST",
			body: formData,
		})
			.then((response) => {
				if (response.redirected) {
					window.location.href = response.url; // Follow redirect on success
				} else {
					alert("Ошибка авторизации");
				}
			})
			.catch((error) => console.error("Error during login:", error));
	});

	// Function to initiate external provider login (redirect to /authorize with provider)
	function initiateExternalLogin(provider) {
		const authorizeUrl = `/authorize?client_id=${encodeURIComponent(clientId)}&redirect_uri=${encodeURIComponent(redirectUri)}&state=${encodeURIComponent(state)}&code_challenge=${encodeURIComponent(codeChallenge)}&nonce=${encodeURIComponent(nonce)}&provider=${encodeURIComponent(provider)}`;
		window.location.href = authorizeUrl;
	}
});
