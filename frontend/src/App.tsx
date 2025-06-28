import { useState } from "react";
import "./App.css";
import MultipleChoice from "./components/MultipleChoice";
import SearchBar from "./components/SearchBar";
import type { question } from "./types/question";

function App() {
  const [questions, setQuestions] = useState<question[] | null>(null);

  return (
    <div className="home-background flex flex-col justify-center bg-black w-full">
      <SearchBar setQuestions={setQuestions} />
      {questions &&
        questions.map((q, index) => (
          <MultipleChoice key={index} question={q} />
        ))}
    </div>
  );
}

export default App;
