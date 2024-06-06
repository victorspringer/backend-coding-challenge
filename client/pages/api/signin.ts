import { NextApiRequest, NextApiResponse } from 'next';
import fetch from 'isomorphic-fetch';
import { SetCookie } from '../../src/cookies';

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
      SetCookie(res, [
        { maxAge: data.response.accessTokenExpiration, key: "MRSAccessToken", value: data.response.accessToken },
        { maxAge: data.response.refreshTokenExpiration, key: "MRSRefreshToken", value: data.response.refreshToken },
        { maxAge: data.response.refreshTokenExpiration, key: "username", value: req.body.username },
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
