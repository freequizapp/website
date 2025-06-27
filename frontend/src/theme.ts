// theme.ts
import { createTheme } from "@mui/material/styles";

const theme = createTheme({
  palette: {
    primary: {
      main: "#3b82f6", // e.g. Tailwind blue-500
    },
    secondary: {
      main: "#a833ff", // Tailwind rose-500
    },
    text: {
      primary: "#ffffff",
    },
    background: {
      default: "#0f172a", // Tailwind slate-900 or similar
    },
  },
  components: {
    MuiRadio: {
      styleOverrides: {
        root: {
          color: "#a388ff",
          "&.Mui-checked": {
            color: "#a833ff", // checked
          },
        },
      },
    },
    MuiFormLabel: {
      styleOverrides: {
        root: {
          color: "#ffffff",
        },
      },
    },
  },
});

export default theme;
