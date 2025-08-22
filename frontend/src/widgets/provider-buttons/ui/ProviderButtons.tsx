import type { Provider } from "../../../entities/provider/types";
import { handleOAuthRedirect } from "../../../features/auth/lib/handleOAuthRedirect";
import LoginForm from "../../../features/login/ui/LoginForm";
import "./ProviderButtons.css";

const ProviderButtons = ({
	providers,
	clientId,
}: {
	providers: Provider[];
	clientId: string;
}) => {
	return (
		<div className="provider-buttons">
			{providers.map((provider) => {
				if (provider === "login_password") {
					return <LoginForm key={provider} clientId={clientId} />;
				}
				return (
					<button
						key={provider}
						data-provider={provider}
						onClick={() => handleOAuthRedirect(provider, clientId)}
					>
						Login with{" "}
						{provider.charAt(0).toUpperCase() +
							provider.slice(1).replace("_", "/")}
					</button>
				);
			})}
		</div>
	);
};

export default ProviderButtons;
