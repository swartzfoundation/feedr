import { useState } from "react";
import { Button } from "./components/ui/button";
import Layout from "./components/layout";

function App() {
  const [count, setCount] = useState(0);

  return (
    <Layout>
      <p>{count}</p>
      <Button onClick={() => setCount((count) => count + 1)}>Increment</Button>
    </Layout>
  );
}

export default App;
