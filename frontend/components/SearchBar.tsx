import { TextField, InputAdornment } from "@mui/material";
import Search from "@mui/icons-material/Search";

function SearchBar() {
  return (
    <div className="my-8">
      <TextField
        label="Generate a quiz about..."
        variant="outlined"
        fullWidth
        InputProps={{
          startAdornment: (
            <InputAdornment position="start">
              <Search />
            </InputAdornment>
          ),
        }}
      />
    </div>
  );
}

export default SearchBar;
