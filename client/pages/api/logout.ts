import { NextApiRequest, NextApiResponse } from 'next';
import fetch from 'isomorphic-fetch';
import { SetCookie } from '../../src/cookies';

export default async (req: NextApiRequest, res: NextApiResponse) => {
  if (req.method === 'POST') {
    const body = JSON.parse(req.body);

    const response = await fetch(`${process.env.NEXT_PUBLIC_AUTH_SERVICE_URL}/logout`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${body.accessToken}`,
      },
    });
    const data = await response.json();

    SetCookie(res, [
      { maxAge: -1, key: "MRSAccessToken", value: "" },
      { maxAge: -1, key: "MRSRefreshToken", value: "" },
      { maxAge: -1, key: "username", value: "" },
    ]);
    
    if (data.statusCode === 200) {
      return res.status(200).json(data);
    } else {
      return res.status(data.statusCode).json(data);
    }
  } else {
    res.setHeader('Allow', ['POST']);
    res.status(405).end(`Method ${req.method} Not Allowed`);
  }
};
