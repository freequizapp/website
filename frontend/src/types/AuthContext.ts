import type { AnswerMap } from "./AnswerMap";
import type { Question } from "./Question";

export type AuthContextType = {
  //user: User | null;
  //setUser: (user: User | null) => void;
  currentQuiz: Question[];
  setCurrentQuiz: React.Dispatch<React.SetStateAction<Question[]>>;
  answers: AnswerMap[];
  setAnswers: React.Dispatch<React.SetStateAction<AnswerMap[]>>;
};
