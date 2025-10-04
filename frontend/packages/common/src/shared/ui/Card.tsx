export const Card: React.FC<{ children: React.ReactNode }> = ({ children }) => (
  <div className="w-full max-w-md p-8 space-y-6 bg-white rounded-lg shadow-md">{children}</div>
);
