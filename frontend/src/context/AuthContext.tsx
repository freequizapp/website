"use client";
import { createContext, useContext, useState } from "react";
import type { ReactNode } from "react";
import type { Question } from "../types/Question";
import type { AuthContextType } from "../types/AuthContext";
import type { AnswerMap } from "../types/AnswerMap";

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: ReactNode }) {
  const [currentQuiz, setCurrentQuiz] = useState<Question[]>([]);
  const [answers, setAnswers] = useState<AnswerMap[]>([]);

  return (
    <AuthContext.Provider
      value={{
        currentQuiz,
        setCurrentQuiz,
        answers,
        setAnswers,
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
