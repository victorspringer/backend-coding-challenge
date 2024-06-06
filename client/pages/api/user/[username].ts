
import { NextApiRequest, NextApiResponse } from 'next';
import fetch from 'isomorphic-fetch';

export default async (req: NextApiRequest, res: NextApiResponse) => {
    if (req.method === 'GET') {
        const response = await fetch(`${process.env.NEXT_PUBLIC_USER_SERVICE_URL}/${req.query.username}`, {
            headers: {
                "Authorization": `Bearer ${req.query.accessToken}`
            },
        });
        const data = await response.json();
        
        if (data.statusCode === 200) {
            return res.status(200).json(data);
        } else {
            return res.status(data.statusCode).json(data);
        }
    } else {
        res.setHeader('Allow', ['GET']);
        res.status(405).end(`Method ${req.method} Not Allowed`);
    }
};
