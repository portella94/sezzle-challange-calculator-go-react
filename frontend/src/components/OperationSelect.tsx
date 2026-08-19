import { MenuItem, TextField } from '@mui/material';
import type { Operation } from '../types';
import { OPERATIONS } from '../config/operations';

interface Props {
  value: Operation;
  onChange: (op: Operation) => void;
  disabled?: boolean;
}

/** Presentational dropdown for choosing the arithmetic operation. */
export function OperationSelect({ value, onChange, disabled }: Props) {
  return (
    <TextField
      select
      label="Operation"
      value={value}
      onChange={(e) => onChange(e.target.value as Operation)}
      disabled={disabled}
      fullWidth
    >
      {OPERATIONS.map((op) => (
        <MenuItem key={op.value} value={op.value}>
          {op.label}
        </MenuItem>
      ))}
    </TextField>
  );
}
