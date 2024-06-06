
import { NextApiRequest, NextApiResponse } from 'next';
import fetch from 'isomorphic-fetch';

export default async (req: NextApiRequest, res: NextApiResponse) => {
    if (req.method === 'POST') {        
        const response = await fetch(`${process.env.NEXT_PUBLIC_RATING_SERVICE_URL}/upsert`, {
            method: 'POST',
            headers: {
                'Authorization': `Bearer ${req.query.accessToken}`,
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                userId: req.body.userId,
                movieId: req.body.movieId,
                value: req.body.value,
            }),
        });
        const data = await response.json();

        if (data.statusCode === 200) {
            return res.status(200).json({ success: true });
        } else {
            return res.status(data.statusCode).json(data);
        }
    } else {
        res.setHeader('Allow', ['POST']);
        res.status(405).end(`Method ${req.method} Not Allowed`);
    }
};
