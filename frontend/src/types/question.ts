export type question = {
  question: string;
  answers: {
    text: string;
    correct: boolean;
    reason: string;
  }[];
};
