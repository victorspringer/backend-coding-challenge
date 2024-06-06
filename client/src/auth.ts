import { IncomingMessage, ServerResponse } from "http";
import { GetCookie, RemoveCookie, SetCookie } from "./cookies";

export default async function IsAuthenticated(req: IncomingMessage & { cookies: Partial<{ [key: string]: string; }> }, res: ServerResponse<IncomingMessage>): Promise<boolean> {
    if (GetCookie(req, "MRSAccessToken")) return true;

    const refreshToken = GetCookie(req, "MRSRefreshToken");
    const username = GetCookie(req, "username");
    if (refreshToken && username) {
        const resp = await refresh(refreshToken, username);
        if (resp.statusCode !== 200) {
            RemoveCookie(res, "MRSRefreshToken");
            RemoveCookie(res, "username");
            return false;
        }

        SetCookie(res, [
            { maxAge: resp.response.accessTokenExpiration, key: "MRSAccessToken", value: resp.response.accessToken },
            { maxAge: resp.response.refreshTokenExpiration, key: "MRSRefreshToken", value: resp.response.refreshToken },
            { maxAge: resp.response.refreshTokenExpiration, key: "username", value: username },
        ]);

        return true;
    }

    return false;
}

const refresh = async (refreshToken: string, username: string): Promise<any> => {
    const response = await fetch(`${process.env.NEXT_PUBLIC_WEBAPP_URL}/api/refreshToken`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({refreshToken, username}),
    });
    const data = await response.json()
    return data;
};
