import { useState } from "react";
import "./App.css";
import MultipleChoice from "./components/MultipleChoice";
import CircularProgress from "@mui/material/CircularProgress";
import Button from "@mui/material/Button";
import SearchBar from "./components/SearchBar";
import type { question } from "./types/question";

function App() {
  const [questions, setQuestions] = useState<question[]>([]);
  const [currentIndex, setCurrentIndex] = useState<number>(0);
  const [isLoading, setIsLoading] = useState<boolean>(false);
  const [answeredMap, setAnsweredMap] = useState<{ [index: number]: boolean }>(
    {},
  );

  const handleNext = () => {
    if (questions && currentIndex < questions.length - 1) {
      setCurrentIndex((prev) => prev + 1);
    }
  };

  const handleAnswerSelected = () => {
    setAnsweredMap((prev) => ({ ...prev, [currentIndex]: true }));
  };

  return (
    <div className="home-background flex flex-col justify-center items-center bg-black w-full h-full">
      {!isLoading && !questions && (
        <SearchBar setQuestions={setQuestions} setIsLoading={setIsLoading} />
      )}
      {isLoading && <CircularProgress color="secondary" />}
      {!isLoading && questions && questions.length > 0 && (
        <>
          <MultipleChoice
            question={questions[currentIndex]}
            onAnswerSelected={handleAnswerSelected}
          />
          <div className="h-10">
            {answeredMap[currentIndex] && (
              <Button
                variant="contained"
                color="secondary"
                onClick={handleNext}
                disabled={currentIndex >= questions.length - 1}
                className="bg-red-500"
              >
                Next
              </Button>
            )}
          </div>
        </>
      )}
    </div>
  );
}

export default App;
