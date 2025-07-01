import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { useAuth } from "../context/AuthContext";
import "../App.css";
import MultipleChoice from "../components/MultipleChoice";
import CircularProgress from "@mui/material/CircularProgress";
import Button from "@mui/material/Button";
import SearchBar from "../components/SearchBar";

export default function Home() {
  const [currentIndex, setCurrentIndex] = useState<number>(0);
  const [isLoading, setIsLoading] = useState<boolean>(false);
  const [answeredMap, setAnsweredMap] = useState<{ [index: number]: boolean }>(
    {},
  );
  const { currentQuiz, answers } = useAuth();
  const navigate = useNavigate();

  const handleNext = () => {
    if (currentQuiz && currentIndex < currentQuiz.length - 1) {
      setCurrentIndex((prev) => prev + 1);
    } else if (currentIndex === 4) {
      navigate("/answers");
    }
    console.log("answered map: ", answeredMap);
  };

  const handleAnswerSelected = () => {
    setAnsweredMap((prev) => ({ ...prev, [currentIndex]: true }));
    console.log("answers", answers);
  };

  return (
    <div className="home-background flex flex-col justify-center items-center bg-black w-full h-full">
      {!isLoading && currentQuiz.length === 0 && (
        <SearchBar setIsLoading={setIsLoading} />
      )}
      {isLoading && <CircularProgress color="secondary" />}
      {!isLoading && currentQuiz && currentQuiz.length > 0 && (
        <>
          <MultipleChoice
            question={currentQuiz[currentIndex]}
            onAnswerSelected={handleAnswerSelected}
          />
          <div className="h-10">
            {currentIndex < 4 ? (
              <Button
                variant="contained"
                color="secondary"
                onClick={handleNext}
                disabled={currentIndex >= currentQuiz.length - 1}
                className="bg-red-500"
              >
                Next
              </Button>
            ) : (
              <Button
                variant="contained"
                color="secondary"
                onClick={handleNext}
                className="bg-red-500"
              >
                See Results
              </Button>
            )}
          </div>
        </>
      )}
    </div>
  );
}
