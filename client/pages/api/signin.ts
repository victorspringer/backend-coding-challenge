import cookie from 'cookie';
import { NextApiRequest, NextApiResponse } from 'next';
import fetch from 'isomorphic-fetch';

export default async (req: NextApiRequest, res: NextApiResponse) => {
  if (req.method === 'POST') {
    const response = await fetch(`${process.env.NEXT_PUBLIC_AUTH_SERVICE_URL}/login`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(req.body),
    });
    const data = await response.json();

    if (data.statusCode === 200) {
      res.setHeader('Set-Cookie', [
        cookie.serialize('MRSAccessToken', data.response.accessToken, {
          httpOnly: true,
          secure: process.env.NODE_ENV === 'production',
          maxAge: data.response.accessTokenExpiration,
          path: '/',
          sameSite: 'lax',
        }),
        cookie.serialize('MRSRefreshToken', data.response.refreshToken, {
          httpOnly: true,
          secure: process.env.NODE_ENV === 'production',
          maxAge: data.response.refreshTokenExpiration,
          path: '/',
          sameSite: 'lax',
        }),
      ]);

      return res.status(200).json(data);
    } else {
      return res.status(data.statusCode).json(data);
    }
  } else {
    res.setHeader('Allow', ['POST']);
    res.status(405).end(`Method ${req.method} Not Allowed`);
  }
};
