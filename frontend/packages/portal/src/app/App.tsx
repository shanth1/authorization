import { useState } from "react";
import { Button } from "@frontend/shared/shared/ui/Button";

function App() {
  const [count, setCount] = useState(0);

  return (
    <>
      <h1>Portal</h1>
      <div>
        <Button>Shared Button</Button>
        <button onClick={() => setCount((count) => count + 1)}>count is {count}</button>
      </div>
    </>
  );
}

export default App;
