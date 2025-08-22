import axios from "axios";
import { apiBaseUrl } from "../config/env";

const api = axios.create({
	baseURL: apiBaseUrl,
	withCredentials: true, // For cookies if needed
});

export default api;
