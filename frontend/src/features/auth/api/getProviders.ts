import api from "../../../shared/lib/api";
import type { Provider } from "../../../entities/provider/types";
import { isDevMode } from "../../../shared/config/env";

export const getProviders = async (clientId: string): Promise<Provider[]> => {
	if (isDevMode) {
		// Mock all providers in dev mode without backend
		return ["login_password", "telegram", "google", "github"];
	}
	const response = await api.get(`/client/${clientId}/providers`);
	return response.data.providers; // Assuming response shape { providers: Provider[] }
};
