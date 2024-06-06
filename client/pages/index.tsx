import * as React from 'react';
import Typography from '@mui/material/Typography';
import { Box } from '@mui/material';
import Popcorn from '../public/svg/popcorn.svg';
import { GetServerSideProps } from 'next';
import cookie from 'cookie';
import RefreshToken from '../src/auth';

export const getServerSideProps: GetServerSideProps = async ({ req, res, query }) => {
  const cookies = req.headers.cookie ? cookie.parse(req.headers.cookie) : null;
  const isAuthenticated = cookies && cookies["MRSAccessToken"];

  if (!isAuthenticated) {
    if (cookies && cookies["MRSRefreshToken"]) {
      let resp = await RefreshToken(cookies["MRSRefreshToken"]);
      if (resp.statusCode !== 200) {
        res.setHeader('Set-Cookie', cookie.serialize("MRSRefreshToken", '', {
          maxAge: -1,
          path: '/',
        }));
        return {
          redirect: {
            destination: '/signin',
            permanent: false,
          },
        };
      }
      res.setHeader('Set-Cookie', [
        cookie.serialize('MRSAccessToken', resp.response.accessToken, {
          httpOnly: true,
          secure: process.env.NODE_ENV === 'production',
          maxAge: resp.response.accessTokenExpiration,
          path: '/',
          sameSite: 'lax',
        }),
        cookie.serialize('MRSRefreshToken', resp.response.refreshToken, {
          httpOnly: true,
          secure: process.env.NODE_ENV === 'production',
          maxAge: resp.response.refreshTokenExpiration,
          path: '/',
          sameSite: 'lax',
        }),
      ]);
    }
  }

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
