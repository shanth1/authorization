import { Logo } from "@common/features/Logo";
import { Card } from "@common/shared/Card";
import { BasicLogin } from "@/features/BasicLogin/BasicLogin";
import { Providers } from "@/features/Providers/Providers";

interface AuthFormProps {
  providers: string[];
}

export const AuthForm: React.FC<AuthFormProps> = ({ providers }) => {
  return (
    <Card>
      <Logo />
      <BasicLogin />

      {providers.length > 0 && (
        <>
          <div className="relative my-4">
            <div className="absolute inset-0 flex items-center">
              <div className="w-full border-t border-gray-300" />
            </div>
            <div className="relative flex justify-center text-sm">
              <span className="px-2 bg-white text-gray-500">Or continue with</span>
            </div>
          </div>
          <Providers providers={providers} />
        </>
      )}
    </Card>
  );
};
