export const handleOAuthRedirect = (provider: string, clientId: string) => {
	// Redirect to backend OAuth endpoint, e.g., /api/auth/google?client_id=xxx
	window.location.href = `/api/auth/${provider}?client_id=${clientId}`;
};
