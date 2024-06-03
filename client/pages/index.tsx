import * as React from 'react';
import Typography from '@mui/material/Typography';

export default function Index() {
  return (
    <Typography className="mrs" variant="h4" component="h1" sx={{ mb: 2 }}>
      <span>Movie</span> <span>Rating</span> <span>System</span>
    </Typography>
  );
}
