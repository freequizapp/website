export type Question = {
  question: string;
  answers: {
    text: string;
    correct: boolean;
    reason: string;
  }[];
};
