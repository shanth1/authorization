import { PageLayout, Loader, ErrorMessage } from "@common/shared/ui";
import { AuthForm } from "@/widgets/AuthForm/AuthForm";
import { useAuthParams } from "./model/hooks/useAuthParams";

export const Login: React.FC = () => {
  const { status, errors, responseData } = useAuthParams();

  return (
    <PageLayout>
      {status === "pending" && <Loader />}
      {status === "error" && <ErrorMessage errors={errors} />}
      {status === "success" && responseData && <AuthForm providers={responseData.providers} />}
    </PageLayout>
  );
};
