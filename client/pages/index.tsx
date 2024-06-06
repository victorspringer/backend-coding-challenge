import * as React from 'react';
import Typography from '@mui/material/Typography';
import { Box } from '@mui/material';
import Popcorn from '../public/svg/popcorn.svg';
import { GetServerSideProps } from 'next';
import IsAuthenticated from '../src/auth';

export const getServerSideProps: GetServerSideProps = async ({ req, res }) => {
  if (!IsAuthenticated(req, res)) return {
    redirect: {
      destination: '/signin',
      permanent: false,
    },
  };

  return { props: {} }
}

export default function Index() {
  return (
    <Box textAlign='center'>
      <Typography className="mrs" variant="h4" component="h1" sx={{ mb: 2 }}>
        <span>Movie</span> <span>Rating</span> <span>System</span>
      </Typography>
      <Popcorn />
    </Box>
  );
}
