export interface AuthParams {
  state: string | null;
  scope: string | null;
  nonce: string | null;
  client_id: string | null;
  redirect_uri: string | null;
  response_type: string | null;
  code_challenge: string | null;
  code_challenge_method: string | null;
}

export type AuthStatus = "pending" | "success" | "error";

export interface AuthState {
  status: AuthStatus;
  errors: string[];
  params: AuthParams | null;
  responseData: { providers: string[] } | null;
}
