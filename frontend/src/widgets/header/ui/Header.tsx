import { useEffect, useState } from "react";
import { Link } from "react-router-dom";

const Header = () => {
	const [isAuthenticated, setIsAuthenticated] = useState(false);

	useEffect(() => {
		// Check localStorage or cookie for token
		const token = localStorage.getItem("token"); // Or read cookie
		setIsAuthenticated(!!token);
	}, []);

	return (
		<header>
			<nav>
				<Link to="/">Home</Link>
				{isAuthenticated ? (
					<>
						<Link to="/profile">Profile</Link>
						<button
							onClick={() => {
								localStorage.removeItem("token");
								setIsAuthenticated(false);
							}}
						>
							Logout
						</button>
					</>
				) : (
					<>
						<Link to="/login">Login</Link>
						<Link to="/registration">Register</Link>
					</>
				)}
			</nav>
		</header>
	);
};

export default Header;
