export type Question = {
  question: string;
  answers: Answer[];
};

export type Answer = {
  text: string;
  correct: boolean;
  reason: string;
};
