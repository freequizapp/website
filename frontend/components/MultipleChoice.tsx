import Radio from "@mui/material/Radio";
import RadioGroup from "@mui/material/RadioGroup";
import FormControlLabel from "@mui/material/FormControlLabel";
import FormControl from "@mui/material/FormControl";
import FormLabel from "@mui/material/FormLabel";

function MultipleChoice() {
  return (
    <div>
      <FormControl className="text-white">
        <FormLabel id="demo-radio-buttons-group-label" className="text-white">
          How do you write a constant variable?
        </FormLabel>
        <RadioGroup
          aria-labelledby="demo-radio-buttons-group-label"
          defaultValue="female"
          name="radio-buttons-group"
          className="text-white"
        >
          <FormControlLabel
            value="female"
            control={<Radio />}
            label="let foo = 0;"
          />
          <FormControlLabel
            value="male"
            control={<Radio />}
            label="const foo = true;"
          />
          <FormControlLabel
            value="other"
            control={<Radio />}
            label="constant foo = 8;"
          />
          <FormControlLabel
            value="other"
            control={<Radio />}
            label="let const foo = 8;"
          />
        </RadioGroup>
      </FormControl>
    </div>
  );
}

export default MultipleChoice;
