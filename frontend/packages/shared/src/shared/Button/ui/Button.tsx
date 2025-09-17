interface ButtonProps {
  children: React.ReactNode;
  variant?: "primary" | "secondary";
}

const Button: React.FC<ButtonProps> = ({ children, variant = "primary" }) => (
  <button
    className={`px-4 py-2 rounded ${variant === "primary" ? "bg-blue-500 text-white" : "bg-gray-200"}`}
  >
    {children}
  </button>
);

export default Button;
