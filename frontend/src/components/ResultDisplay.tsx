import { Alert, Box, CircularProgress, Paper, Typography } from '@mui/material';
import { formatResult } from '../utils/format';

interface Props {
  loading: boolean;
  error: string | null;
  result: number | null;
}

/**
 * Presentational output area. Shows exactly one of: a spinner while loading, an
 * error alert, or the formatted result. Renders nothing before the first
 * calculation.
 */
export function ResultDisplay({ loading, error, result }: Props) {
  if (loading) {
    return (
      <Box sx={{ display: 'flex', justifyContent: 'center', py: 2 }}>
        <CircularProgress aria-label="Calculating" />
      </Box>
    );
  }

  if (error) {
    return (
      <Alert severity="error" role="alert">
        {error}
      </Alert>
    );
  }

  if (result !== null) {
    return (
      <Paper variant="outlined" sx={{ p: 2, textAlign: 'center', bgcolor: 'action.hover' }}>
        <Typography variant="overline" color="text.secondary">
          Result
        </Typography>
        <Typography variant="h4" component="output" data-testid="result" sx={{ wordBreak: 'break-all' }}>
          {formatResult(result)}
        </Typography>
      </Paper>
    );
  }

  return null;
}
