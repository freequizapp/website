import { useEffect } from "react";
import { useAuth } from "../context/AuthContext";
import { useNavigate } from "react-router-dom";
import WrongAnswer from "../components/WrongAnswer";
import { Button } from "@mui/material";
import type { Question } from "../types/Question";

export default function Answers() {
  const { currentQuiz, answers, resetQuiz, setNumberRight, numberRight } =
    useAuth();
  const navigate = useNavigate();

  useEffect(() => {
    let count = 0;

    currentQuiz.forEach((q: Question) => {
      const userAnswerText = answers[q.question];
      const correctAnswer = q.answers.find((a) => a.correct)?.text;

      if (userAnswerText === correctAnswer) {
        count += 1;
      }
    });

    setNumberRight(count);
  }, []);

  if (Object.keys(answers).length > 0) {
    return (
      <div>
        <h2 className="font-black text-2xl">Results</h2>
        <h3 className="font-black text-xl">{numberRight}/5</h3>
        {currentQuiz.map((q: Question, index) => {
          const selectedAnswerText: string = answers[q.question];

          const isCorrect = q.answers.find(
            (a) => a.text === selectedAnswerText && a.correct,
          );

          return !isCorrect ? (
            <WrongAnswer question={q} answer={selectedAnswerText} key={index} />
          ) : null;
        })}
        <div>
          <Button
            variant="contained"
            color="secondary"
            onClick={() => {
              navigate("/");
              resetQuiz();
            }}
          >
            Generate New Quiz
          </Button>
        </div>
      </div>
    );
  } else {
    return (
      <div>
        <Button
          variant="contained"
          color="secondary"
          onClick={() => navigate("/")}
        >
          Generate Quiz
        </Button>
      </div>
    );
  }
}
