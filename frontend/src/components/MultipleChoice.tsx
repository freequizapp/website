import Radio from "@mui/material/Radio";
import RadioGroup from "@mui/material/RadioGroup";
import FormControlLabel from "@mui/material/FormControlLabel";
import FormControl from "@mui/material/FormControl";
import FormLabel from "@mui/material/FormLabel";
import type { question } from "../types/question";

type MultipleChoiceProp = {
  question: question;
  onAnswerSelected: () => void;
};

function MultipleChoice({ question, onAnswerSelected }: MultipleChoiceProp) {
  const handleChange = () => {
    onAnswerSelected();
  };

  return (
    <div className="my-8 w-full min-w-l flex justify-start items-center">
      <FormControl className="text-white">
        <div className="mb-8">
          <FormLabel id="demo-radio-buttons-group-label">
            {question.question}
          </FormLabel>
        </div>
        <RadioGroup
          aria-labelledby="demo-radio-buttons-group-label"
          name="radio-buttons-group"
          onChange={handleChange}
        >
          {question &&
            question.answers.map((answer, index) => (
              <FormControlLabel
                control={<Radio />}
                label={answer.text}
                value={answer.text}
                key={index}
              />
            ))}
        </RadioGroup>
      </FormControl>
    </div>
  );
}

export default MultipleChoice;
