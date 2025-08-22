import { useState } from "react";
import api from "../../../shared/lib/api";

const LoginForm = ({ clientId }: { clientId?: string }) => {
	// Optional clientId for authorize context
	const [login, setLogin] = useState("");
	const [password, setPassword] = useState("");
	const [error, setError] = useState("");

	const handleSubmit = async (e: React.FormEvent) => {
		e.preventDefault();
		try {
			const response = await api.post("/login", {
				login,
				password,
				client_id: clientId,
			});
			// Assume sets cookie or localStorage token
			localStorage.setItem("token", response.data.token); // Or use cookies
			window.location.href = "/profile"; // Or back to client redirect_uri if OAuth
		} catch (err) {
			setError("Login failed");
		}
	};

	return (
		<form onSubmit={handleSubmit}>
			<input
				type="text"
				placeholder="Login"
				value={login}
				onChange={(e) => setLogin(e.target.value)}
				required
			/>
			<input
				type="password"
				placeholder="Password"
				value={password}
				onChange={(e) => setPassword(e.target.value)}
				required
			/>
			{error && <p>{error}</p>}
			<button type="submit">Login</button>
		</form>
	);
};

export default LoginForm;
