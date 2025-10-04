import { AuthParams } from "@/entities/auth/model/types";
import { REQUIRED_SCOPES, VALID_CODE_CHALLENGE_METHOD } from "../config/auth";

VALID_CODE_CHALLENGE_METHOD;
export const getValidationErrors = (params: AuthParams): string[] => {
  const errors: string[] = [];

  for (const key in params) {
    if (!params[key as keyof AuthParams]) {
      errors.push(`Missing required parameter: ${key}.`);
    }
  }

  if (errors.length > 0) return errors;

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
