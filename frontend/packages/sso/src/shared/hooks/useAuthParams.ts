import { useState, useEffect } from "react";
import { useLocation } from "react-router-dom";
import { AuthParams, AuthState } from "@/entities/auth/model/types";
import { getValidationErrors } from "./getValidationErrors";

export const useAuthParams = (): AuthState => {
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

      try {
        // TODO: Replace with real API request
        const mockServerResponse: { providers: string[] } = await new Promise((resolve) =>
          setTimeout(() => resolve({ providers: ["Google", "GitHub", "Telegram"] }), 1000)
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

  return authState;
};
