import { useState } from "react";
import { useAuth } from "../context/AuthContext";
import api_url from "../config/config";
import type { Dispatch, SetStateAction } from "react";
import type { Question } from "../types/Question";
import { TextField, InputAdornment } from "@mui/material";
import Search from "@mui/icons-material/Search";

type Props = {
  setIsLoading: Dispatch<SetStateAction<boolean>>;
};

function SearchBar({ setIsLoading }: Props) {
  const { setCurrentQuiz } = useAuth();
  const [searchText, setSearchText] = useState<string>("");

  const handleKeyDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
    if (e.key === "Enter") {
      e.preventDefault(); // prevent form submission if inside a form
      generateSearch();
    }
  };

  const generateSearch = async () => {
    setCurrentQuiz([]);
    setIsLoading(true);

    const res = await fetch(api_url, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        prompt: searchText,
      }),
    });

    if (!res.body) throw new Error("no response body");

    const reader = res.body.getReader();
    const decoder = new TextDecoder();
    let buffer = "";

    const processBuffer = () => {
      let newLineIndex;
      while ((newLineIndex = buffer.indexOf("\n")) !== -1) {
        const line = buffer.slice(0, newLineIndex).trim();
        buffer = buffer.slice(newLineIndex + 1);
        if (!line) continue;
        try {
          const question = JSON.parse(line);
          setCurrentQuiz((prev: Question[]) => [...(prev ?? []), question]);
          setIsLoading(false);
        } catch (error) {
          console.error("failed to parse streamed chunk: ", line);
        }
      }
    };

    while (true) {
      const { done, value } = await reader.read();
      if (done) break;
      buffer += decoder.decode(value, { stream: true });
      processBuffer();
    }

    // final flush, incase last chunk had no trailing \n
    buffer += decoder.decode();
    processBuffer();
  };

  return (
    <div className="my-8">
      <TextField
        label="Generate a quiz about..."
        variant="outlined"
        fullWidth
        value={searchText}
        onChange={(e) => setSearchText(e.target.value)}
        onKeyDown={handleKeyDown}
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
