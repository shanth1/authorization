import { Loader } from "@common/shared/Loader";
import { ErrorMessage } from "@common/shared/ErrorMessage";
import { PageLayout } from "@common/shared/PageLayout";
import { useAuthParams } from "@/shared/hooks/useAuthParams";
import { AuthForm } from "@/widgets/AuthForm/AuthForm";

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
