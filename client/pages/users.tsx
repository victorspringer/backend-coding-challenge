import * as React from 'react';
import Typography from '@mui/material/Typography';
import Button from '@mui/material/Button';
import Link from '../src/Link';
import Box from '@mui/material/Box';

export default function Users() {
  return (
    <>
      <Typography variant="h4" component="h1" sx={{ mb: 2 }}>
        Users page
      </Typography>
      <Box sx={{ maxWidth: 'sm' }}>
        <Button variant="contained" component={Link} noLinkStyle href="/">
          Go to the home page
        </Button>
      </Box>
    </>
  );
}
