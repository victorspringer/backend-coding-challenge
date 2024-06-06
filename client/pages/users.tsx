import * as React from 'react';
import Typography from '@mui/material/Typography';
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

export default function Users() {
  return (
    <Typography variant="h4" component="h1" sx={{ mb: 2 }}>
      Users page
    </Typography>
  );
}
