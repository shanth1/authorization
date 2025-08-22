import { useEffect } from "react";
import { useNavigate, useSearchParams } from "react-router-dom";
import { useQuery } from "@tanstack/react-query";
import { getProviders } from "../../../features/auth/api/getProviders";
import Loader from "../../../widgets/loader/ui/Loader.tsx";
import ProviderButtons from "../../../widgets/provider-buttons/ui/ProviderButtons.tsx";

const AuthorizePage = () => {
	const [searchParams] = useSearchParams();
	const navigate = useNavigate();
	const clientId = searchParams.get("client_id");

	useEffect(() => {
		if (!clientId) {
			navigate("/"); // Redirect if no client_id
		}
	}, [clientId, navigate]);

	const { data: providers, isLoading } = useQuery({
		queryKey: ["providers", clientId],
		queryFn: () => getProviders(clientId!),
		enabled: !!clientId,
	});

	if (!clientId) return null; // Handled by useEffect

	if (isLoading) return <Loader />;

	return (
		<div>
			<h1>Authorize</h1>
			<ProviderButtons providers={providers || []} clientId={clientId} />
			<a href="/registration">Register new account</a>{" "}
			{/* Link to registration from authorize */}
		</div>
	);
};

export default AuthorizePage;
