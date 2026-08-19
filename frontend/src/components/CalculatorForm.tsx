import { Box, Button, Stack, TextField, Typography } from '@mui/material';
import type { FormEvent } from 'react';
import { useCalculator } from '../hooks/useCalculator';
import { OperationSelect } from './OperationSelect';
import { ResultDisplay } from './ResultDisplay';

/**
 * Container that wires the calculator state hook to the presentational inputs
 * and output. Renders one operand input per operand the selected operation
 * requires (driven by the operation metadata's arity).
 */
export function CalculatorForm() {
  const { operation, meta, operands, result, error, fieldErrors, loading, setOperation, setOperand, submit } =
    useCalculator();

  const handleSubmit = (e: FormEvent) => {
    e.preventDefault();
    void submit();
  };

  return (
    <Box component="form" onSubmit={handleSubmit} noValidate aria-label="Calculator">
      <Stack spacing={2.5}>
        <OperationSelect value={operation} onChange={setOperation} disabled={loading} />

        <Typography variant="body2" color="text.secondary">
          {meta.hint}
        </Typography>

        {Array.from({ length: meta.arity }).map((_, index) => (
          <TextField
            key={index}
            label={meta.operandLabels[index]}
            value={operands[index] ?? ''}
            onChange={(e) => setOperand(index, e.target.value)}
            error={Boolean(fieldErrors[index])}
            helperText={fieldErrors[index] ?? ' '}
            disabled={loading}
            fullWidth
            inputMode="decimal"
            slotProps={{ htmlInput: { inputMode: 'decimal' } }}
          />
        ))}

        <Button type="submit" variant="contained" size="large" disabled={loading}>
          Calculate
        </Button>

        <ResultDisplay loading={loading} error={error} result={result} />
      </Stack>
    </Box>
  );
}
