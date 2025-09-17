import { useState, useEffect } from "react";
import { useLocation } from "react-router-dom";

interface AuthParams {
  state: string | null;
  scope: string | null;
  nonce: string | null;
  code_challenge: string | null;
  code_challenge_method: string | null;
}

const LoginPage: React.FC = () => {
  const location = useLocation();
  const [params, setParams] = useState<AuthParams>({
    state: null,
    scope: null,
    nonce: null,
    code_challenge: null,
    code_challenge_method: null,
  });
  const [isValid, setIsValid] = useState<boolean>(false);
  const [isChecking, setIsChecking] = useState<boolean>(true);

  useEffect(() => {
    const searchParams = new URLSearchParams(location.search);

    const extractedParams: AuthParams = {
      state: searchParams.get("state"),
      scope: searchParams.get("scope"),
      nonce: searchParams.get("nonce"),
      code_challenge: searchParams.get("code_challenge"),
      code_challenge_method: searchParams.get("code_challenge_method"),
    };

    setParams(extractedParams);
    validateParams(extractedParams);
  }, [location]);

  const validateParams = (params: AuthParams) => {
    const { state, scope, nonce, code_challenge, code_challenge_method } = params;

    const hasRequiredParams = Boolean(
      state && scope && nonce && code_challenge && code_challenge_method
    );

    const isCodeChallengeMethodValid = code_challenge_method === "S256";
    const isScopeValid =
      scope?.includes("openid") && scope?.includes("profile") && scope?.includes("email");

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
