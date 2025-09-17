import { LoginPage } from "../pages";
import { BrowserRouter as Router, Routes, Route } from "react-router-dom";

const App: React.FC = () => {
  return (
    <Router>
      <Routes>
        <Route path="*" element={<LoginPage />} />
      </Routes>
    </Router>
  );
};

export default App;
