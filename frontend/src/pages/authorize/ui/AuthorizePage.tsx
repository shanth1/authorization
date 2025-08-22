import { useEffect } from "react";
import { useNavigate, useSearchParams } from "react-router-dom";
import { useQuery } from "@tanstack/react-query";
import { getProviders } from "../../../features/auth/api/getProviders";
import Loader from "../../../widgets/loader/ui/Loader";
import ProviderButtons from "../../../widgets/provider-buttons/ui/ProviderButtons";

const AuthorizePage = () => {
	const [searchParams] = useSearchParams();
	const navigate = useNavigate();
	const clientId = searchParams.get("client_id");

	useEffect(() => {
		if (!clientId) {
			navigate("/");
		}
	}, [clientId, navigate]);

	const { data: providers, isLoading } = useQuery({
		queryKey: ["providers", clientId],
		queryFn: () => getProviders(clientId!),
		enabled: !!clientId,
	});

	if (!clientId) return null;

	if (isLoading) return <Loader />;

	return (
		<div className="container">
			<h1>Authorize</h1>
			<ProviderButtons providers={providers || []} clientId={clientId} />
			<a href="/registration">Register new account</a>
		</div>
	);
};

export default AuthorizePage;
