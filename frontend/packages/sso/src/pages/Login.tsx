import { Loader } from "@frontend/shared/shared/Loader";
import { ErrorMessage } from "@frontend/shared/shared/ErrorMessage";
import { PageLayout } from "@frontend/shared/shared/PageLayout";
import { useAuthParams } from "../shared/hooks/useAuthParams";
import { AuthForm } from "../widgets/AuthForm/AuthForm";

export const LoginPage: React.FC = () => {
  const { status, errors, responseData } = useAuthParams();

  return (
    <PageLayout>
      {status === "pending" && <Loader />}
      {status === "error" && <ErrorMessage errors={errors} />}
      {status === "success" && responseData && <AuthForm providers={responseData.providers} />}
    </PageLayout>
  );
};
