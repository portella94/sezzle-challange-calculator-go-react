import { describe, expect, it } from 'vitest';
import { render, screen } from '@testing-library/react';
import { ThemeProvider } from '@mui/material';
import App from './App';
import { theme } from './theme';

describe('App', () => {
  it('renders the calculator heading and form', () => {
    render(
      <ThemeProvider theme={theme}>
        <App />
      </ThemeProvider>,
    );

    expect(screen.getByRole('heading', { name: 'Calculator', level: 1 })).toBeInTheDocument();
    expect(screen.getByRole('button', { name: /calculate/i })).toBeInTheDocument();
  });
});
