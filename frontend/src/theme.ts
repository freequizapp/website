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
          alignSelf: "flex-start", // aligns radio icon with top of multiline text
          padding: "4px", // reduce default padding to help alignment
          marginRight: "5px",

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
          "&.Mui-focused": {
            color: "#ffffff",
          },
          width: "100%",
          textAlign: "left",
        },
      },
    },
    MuiFormControlLabel: {
      styleOverrides: {
        root: {
          textAlign: "left",
          margin: "5px 0",
        },
      },
    },
    MuiOutlinedInput: {
      styleOverrides: {
        root: {
          "& fieldset": {
            borderColor: "#ffffff", // default (not focused)
          },
        },
      },
    },
  },
});

export default theme;
