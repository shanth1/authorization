import { useState, useEffect } from "react";
import { useLocation } from "react-router-dom";

interface AuthParams {
  state: string | null;
  scope: string | null;
  nonce: string | null;
  clientId: string | null;
  redirectUri: string | null;
  responseType: string | null;
  codeChallenge: string | null;
  codeChallengeMethod: string | null;
}

const LoginPage: React.FC = () => {
  const location = useLocation();
  const [_, setParams] = useState<AuthParams>();
  const [isValid, setIsValid] = useState<boolean>(false);
  const [isChecking, setIsChecking] = useState<boolean>(true);

  useEffect(() => {
    const searchParams = new URLSearchParams(location.search);

    const extractedParams: AuthParams = {
      state: searchParams.get("state"),
      scope: searchParams.get("scope"),
      nonce: searchParams.get("nonce"),
      clientId: searchParams.get("client_id"),
      redirectUri: searchParams.get("redirect_uri"),
      responseType: searchParams.get("response_type"),
      codeChallenge: searchParams.get("code_challenge"),
      codeChallengeMethod: searchParams.get("code_challenge_method"),
    };

    setParams(extractedParams);
    validateParams(extractedParams);
  }, [location]);

  const validateParams = (params: AuthParams) => {
    const hasRequiredParams = Object.values(params).every(
      (val) => val !== undefined && val !== null && val !== ""
    );

    const isCodeChallengeMethodValid = params.codeChallengeMethod === "S256";
    const isScopeValid =
      params?.scope?.includes("openid") &&
      params?.scope?.includes("profile") &&
      params.scope?.includes("email");

    const isValid = hasRequiredParams && isCodeChallengeMethodValid && isScopeValid;

    setIsValid(!!isValid);
    setIsChecking(false);
  };

  if (isChecking) {
    return <div>Loading</div>;
  }

  return isValid ? <div>Yup</div> : <div>Nope</div>;
};

export default LoginPage;
