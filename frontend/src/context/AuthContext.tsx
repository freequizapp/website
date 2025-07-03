"use client";
import { createContext, useContext, useState } from "react";
import type { ReactNode } from "react";
import type { Question } from "../types/Question";

type AuthContextType = {
  //user: User | null;
  //setUser: (user: User | null) => void;
  currentQuiz: Question[];
  setCurrentQuiz: React.Dispatch<React.SetStateAction<Question[]>>;
  answers: Record<string, string>;
  setAnswers: React.Dispatch<React.SetStateAction<Record<string, string>>>;
  resetQuiz: () => void;
  numberRight: number;
  setNumberRight: React.Dispatch<React.SetStateAction<number>>;
};

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: ReactNode }) {
  const [currentQuiz, setCurrentQuiz] = useState<Question[]>([]);
  const [answers, setAnswers] = useState<Record<string, string>>({});
  const [numberRight, setNumberRight] = useState<number>(0);

  const resetQuiz = () => {
    console.log("resetting state for quiz and answers");
    setCurrentQuiz([]);
    setAnswers({});
    setNumberRight(0);
  };

  return (
    <AuthContext.Provider
      value={{
        currentQuiz,
        setCurrentQuiz,
        answers,
        setAnswers,
        resetQuiz,
        numberRight,
        setNumberRight,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error("useAuth must be used within an AuthProvider");
  }
  return context;
}
