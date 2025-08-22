export const isDevMode =
	import.meta.env.DEV || import.meta.env.VITE_DEV_MODE === "true";
export const apiBaseUrl = "/api"; // Proxied in dev, direct in prod
