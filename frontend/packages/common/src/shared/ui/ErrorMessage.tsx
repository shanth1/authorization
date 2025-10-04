export const ErrorMessage: React.FC<{ errors: string[] }> = ({ errors }) => (
  <div className="p-4 border border-red-300 bg-red-50 rounded-lg">
    <h2 className="text-lg font-semibold text-red-800">Authorization Error</h2>
    <p className="mt-2 text-sm text-red-700">Unable to process the request due to:</p>
    <ul className="mt-2 list-disc list-inside space-y-1 text-sm text-red-700">
      {errors.map((error, index) => (
        <li key={index}>{error}</li>
      ))}
    </ul>
  </div>
);
