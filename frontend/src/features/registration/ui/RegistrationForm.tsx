import { useState } from "react";
import api from "../../../shared/lib/api";
// Add form validation library like react-hook-form if needed for best practices

const RegistrationForm = () => {
	const [name, setName] = useState("");
	const [login, setLogin] = useState("");
	const [email, setEmail] = useState("");
	const [password, setPassword] = useState("");
	const [confirmPassword, setConfirmPassword] = useState("");
	const [error, setError] = useState("");

	const handleSubmit = async (e: React.FormEvent) => {
		e.preventDefault();
		if (password !== confirmPassword) {
			setError("Passwords do not match");
			return;
		}
		try {
			await api.post("/register", { name, login, email, password });
			// Redirect to /login or home after success
			window.location.href = "/login";
		} catch (err) {
			setError("Registration failed");
		}
	};

	return (
		<form onSubmit={handleSubmit}>
			<input
				type="text"
				placeholder="Name"
				value={name}
				onChange={(e) => setName(e.target.value)}
				required
			/>
			<input
				type="text"
				placeholder="Login"
				value={login}
				onChange={(e) => setLogin(e.target.value)}
				required
			/>
			<input
				type="email"
				placeholder="Email"
				value={email}
				onChange={(e) => setEmail(e.target.value)}
				required
			/>
			<input
				type="password"
				placeholder="Password"
				value={password}
				onChange={(e) => setPassword(e.target.value)}
				required
			/>
			<input
				type="password"
				placeholder="Confirm Password"
				value={confirmPassword}
				onChange={(e) => setConfirmPassword(e.target.value)}
				required
			/>
			{error && <p>{error}</p>}
			<button type="submit">Register</button>
		</form>
	);
};

export default RegistrationForm;
