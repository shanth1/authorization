import { useState, useEffect } from "react";
import { useLocation } from "react-router-dom";

interface AuthParams {
  state: string | null;
  scope: string | null;
  nonce: string | null;
  client_id: string | null;
  redirect_uri: string | null;
  response_type: string | null;
  code_challenge: string | null;
  code_challenge_method: string | null;
}

type AuthStatus = "pending" | "success" | "error";

interface AuthState {
  status: AuthStatus;
  errors: string[];
  params: AuthParams | null;
  responseData: any | null;
}

const REQUIRED_SCOPES = ["openid", "profile", "email"];
const VALID_CODE_CHALLENGE_METHOD = "S256";

const getValidationErrors = (params: AuthParams): string[] => {
  const errors: string[] = [];

  for (const key in params) {
    const value = params[key as keyof AuthParams];
    if (!value) {
      errors.push(`Missing required parameter: ${key}.`);
    }
  }

  if (errors.length > 0) {
    return errors;
  }

  if (params.code_challenge_method !== VALID_CODE_CHALLENGE_METHOD) {
    errors.push(`Invalid code_challenge_method. Expected "${VALID_CODE_CHALLENGE_METHOD}".`);
  }

  const scopes = params.scope?.split(" ") || [];
  REQUIRED_SCOPES.forEach((requiredScope) => {
    if (!scopes.includes(requiredScope)) {
      errors.push(`Missing required value in scope parameter: "${requiredScope}".`);
    }
  });

  return errors;
};

const LoginPage: React.FC = () => {
  const location = useLocation();
  const [authState, setAuthState] = useState<AuthState>({
    status: "pending",
    errors: [],
    params: null,
    responseData: null,
  });

  useEffect(() => {
    const processAuthRequest = async () => {
      const searchParams = new URLSearchParams(location.search);

      const extractedParams: AuthParams = {
        state: searchParams.get("state"),
        scope: searchParams.get("scope"),
        nonce: searchParams.get("nonce"),
        client_id: searchParams.get("client_id"),
        redirect_uri: searchParams.get("redirect_uri"),
        response_type: searchParams.get("response_type"),
        code_challenge: searchParams.get("code_challenge"),
        code_challenge_method: searchParams.get("code_challenge_method"),
      };

      const validationErrors = getValidationErrors(extractedParams);
      if (validationErrors.length > 0) {
        setAuthState({
          status: "error",
          errors: validationErrors,
          params: extractedParams,
          responseData: null,
        });
        return;
      }

      // TODO: Replace with real API request
      try {
        const mockServerResponse = await new Promise((resolve) =>
          setTimeout(
            () =>
              resolve({
                providers: ["Telegram", "Google"],
              }),
            1000
          )
        );

        setAuthState({
          status: "success",
          errors: [],
          params: extractedParams,
          responseData: mockServerResponse,
        });
      } catch (error) {
        const errorMessage = error instanceof Error ? error.message : "An unknown error occurred.";
        setAuthState({
          status: "error",
          errors: [errorMessage],
          params: extractedParams,
          responseData: null,
        });
      }
    };

    processAuthRequest();
  }, [location]);

  if (authState.status === "pending") {
    return <div>Validating authorization parameters...</div>;
  }

  if (authState.status === "error") {
    return (
      <div>
        <h1>Authorization Error</h1>
        <p>Unable to process the request due to the following reasons:</p>
        <ul>
          {authState.errors.map((error, index) => (
            <li key={index}>{error}</li>
          ))}
        </ul>
      </div>
    );
  }

  if (authState.status === "success") {
    return (
      <div>
        <h1>Choose a login method</h1>
        {authState.responseData?.providers?.map((provider: string) => (
          <button key={provider} onClick={() => alert(`Login via ${provider} selected`)}>
            Login with {provider}
          </button>
        ))}
      </div>
    );
  }

  return null;
};

export default LoginPage;
