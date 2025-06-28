import { useState } from "react";
import type { Dispatch, SetStateAction } from "react";
import { TextField, InputAdornment } from "@mui/material";
import Search from "@mui/icons-material/Search";
import type { question } from "../types/question";

type Props = {
  setQuestions: Dispatch<SetStateAction<question[]>>;
};

function SearchBar({ setQuestions }: Props) {
  const [searchText, setSearchText] = useState<string>("");

  const generateSearch = async () => {
    try {
      console.log("search text: ", searchText);
      const res = await fetch("http://localhost:8080/generate-questions", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          prompt: searchText,
        }),
      });

      const newQuestions: question[] = await res.json();
      setQuestions(newQuestions);
      console.log(newQuestions);
    } catch (err) {
      console.error(err);
    }
  };

  return (
    <div className="my-8">
      <TextField
        label="Generate a quiz about..."
        variant="outlined"
        fullWidth
        value={searchText}
        onChange={(e) => setSearchText(e.target.value)}
        InputProps={{
          startAdornment: (
            <InputAdornment position="start">
              <Search className="text-white" onClick={generateSearch} />
            </InputAdornment>
          ),
        }}
      />
    </div>
  );
}

export default SearchBar;
