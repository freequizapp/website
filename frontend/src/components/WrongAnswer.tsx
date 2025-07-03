import type { Question, Answer } from "../types/Question";
import { useEffect, useState } from "react";

type WrongAnswerProp = {
  question: Question;
  answer: string;
};

function WrongAnswer({ question, answer }: WrongAnswerProp) {
  const [wrongAnswer, setWrongAnswer] = useState<Answer>();
  const [correctAnswer, setCorrectAnswer] = useState<Answer>();

  useEffect(() => {
    const wrong = question.answers.filter((a) => a.text === answer);
    setWrongAnswer(wrong[0]);
    const correct = question.answers.filter((a) => a.correct === true);
    setCorrectAnswer(correct[0]);
  }, []);

  return (
    <div className="my-8 w-full min-w-l flex flex-col justify-start items-center my-5">
      <p className="font-semibold">{question.question}</p>
      <p className="melon">Your answer: {wrongAnswer?.text}</p>
      <p className="green">Correct answer: {correctAnswer?.text}</p>
      <p>{wrongAnswer?.reason}</p>
      <p>{correctAnswer?.reason}</p>
    </div>
  );
}

export default WrongAnswer;
