interface ProviderLoginButtonsProps {
  providers: string[];
}

export const Providers: React.FC<ProviderLoginButtonsProps> = ({ providers }) => {
  const handleProviderClick = (provider: string) => {
    alert(`Login via ${provider} selected`);
    // TODO: Implement provider login logic
  };

  return (
    <div className="space-y-3">
      {providers.map((provider) => (
        <button
          key={provider}
          onClick={() => handleProviderClick(provider)}
          className="block bg-gray-100 hover:bg-gray-200"
        >
          Continue with {provider}
        </button>
      ))}
    </div>
  );
};
