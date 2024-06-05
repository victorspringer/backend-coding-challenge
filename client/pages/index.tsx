import * as React from 'react';
import Typography from '@mui/material/Typography';
import { Box } from '@mui/material';
import Popcorn from '../public/svg/popcorn.svg';
import { GetServerSideProps } from 'next';
import cookie from 'cookie';

export const getServerSideProps: GetServerSideProps = async ({ req, res, query }) => {
  const cookies = req.headers.cookie ? cookie.parse(req.headers.cookie) : null;
  const isAuthenticated = cookies && cookies["MRSAccessToken"];

  if (!isAuthenticated) {
      return {
          redirect: {
              destination: '/signin',
              permanent: false,
          },
      };
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
