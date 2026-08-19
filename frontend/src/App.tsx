import { Card, CardContent, Container, Stack, Typography } from '@mui/material';
import { CalculatorForm } from './components/CalculatorForm';

/** Top-level layout: a centered, responsive card holding the calculator. */
export default function App() {
  return (
    <Container maxWidth="sm" sx={{ py: { xs: 3, sm: 6 } }}>
      <Stack spacing={3}>
        <Stack spacing={0.5}>
          <Typography variant="h1" component="h1">
            Calculator
          </Typography>
          <Typography variant="body1" color="text.secondary">
            Basic and advanced arithmetic, powered by a Go REST API.
          </Typography>
        </Stack>

        <Card elevation={3}>
          <CardContent sx={{ p: { xs: 2, sm: 3 } }}>
            <CalculatorForm />
          </CardContent>
        </Card>
      </Stack>
    </Container>
  );
}
