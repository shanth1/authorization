import { Routes, Route } from "react-router-dom";
import Header from "../widgets/header/ui/Header.tsx";
import HomePage from "../pages/home/ui/HomePage.tsx";
import AuthorizePage from "../pages/authorize/ui/AuthorizePage.tsx";
import LoginPage from "../pages/login/ui/LoginPage.tsx";
import RegistrationPage from "../pages/registration/ui/RegistrationPage.tsx";
import ProfilePage from "../pages/profile/ui/ProfilePage.tsx";

const App = () => {
	return (
		<div className="app">
			<Header />
			<main>
				<Routes>
					<Route path="/" element={<HomePage />} />
					<Route path="/authorize" element={<AuthorizePage />} />
					<Route path="/login" element={<LoginPage />} />
					<Route
						path="/registration"
						element={<RegistrationPage />}
					/>
					<Route path="/profile" element={<ProfilePage />} />
					{/* Add 404 or other routes as needed */}
				</Routes>
			</main>
		</div>
	);
};

export default App;
