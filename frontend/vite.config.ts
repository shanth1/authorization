import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";

// https://vitejs.dev/config/
export default defineConfig(({ mode }) => {
	return {
		plugins: [react()],
		base: mode === "production" ? "/" : "/", // Preserve paths in build
		server: {
			proxy: {
				"/api": {
					target: "http://localhost:8080", // Assuming backend runs on 8080; adjust as needed
					changeOrigin: true,
					secure: false,
				},
			},
		},
		build: {
			outDir: "../backend/static", // Build directly into backend's static folder; adjust path if needed
			emptyOutDir: true,
		},
	};
});
