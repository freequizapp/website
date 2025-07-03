import { Routes, Route } from "react-router-dom";
import "./App.css";
import Home from "./pages/Home";
import Answers from "./pages/Answers";

function App() {
  return (
    <Routes>
      <Route path="/" element={<Home />} />
      <Route path="/answers" element={<Answers />} />
    </Routes>
  );
}

export default App;
