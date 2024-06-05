import * as React from 'react';
import Typography from '@mui/material/Typography';
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

export default function Users() {
  return (
    <Typography variant="h4" component="h1" sx={{ mb: 2 }}>
      Users page
    </Typography>
  );
}
