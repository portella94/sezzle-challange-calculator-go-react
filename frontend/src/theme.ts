import { createTheme } from '@mui/material/styles';

/** Application theme. Kept minimal — colors and rounded shapes only. */
export const theme = createTheme({
  palette: {
    mode: 'light',
    primary: { main: '#5b3df5' },
    background: { default: '#f4f5fb' },
  },
  shape: { borderRadius: 12 },
  typography: {
    h1: { fontSize: '2rem', fontWeight: 700 },
  },
});
