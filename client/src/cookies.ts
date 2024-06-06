import { serialize, parse } from "cookie";
import { IncomingMessage, ServerResponse } from "http";

type cookieProps = {
  maxAge: number;
  key: string;
  value: string
};

export function GetCookie(req: IncomingMessage & { cookies: Partial<{ [key: string]: string; }> }, key: string): string | undefined {
  const cookies = parseCookies(req);
  return cookies[key];
}

export function SetCookie(res: ServerResponse<IncomingMessage>, props: cookieProps[]): void {
  const cookies = props.map((p) => {
    return serialize(p.key, p.value, {
      maxAge: p.maxAge,
      expires: new Date(Date.now() + p.maxAge * 1000),
      httpOnly: true,
      secure: process.env.NODE_ENV === "production",
      path: "/",
      sameSite: "lax",
    });
  });

  res.setHeader("Set-Cookie", cookies);
}

export function RemoveCookie(res: ServerResponse<IncomingMessage>, key: string): void {
  const cookie = serialize(key, "", {
    maxAge: -1,
    path: "/",
  });
  
  res.setHeader("Set-Cookie", cookie);
}

export function parseCookies(req: IncomingMessage & { cookies: Partial<{ [key: string]: string; }> }): { [key: string]: string } {
  // for API Routes we don't need to parse the cookies.
  if (req.cookies) return req.cookies as { [key: string]: string };

  // for pages we do need to parse the cookies.
  const cookie = req.headers?.cookie;

  return parse(cookie || "");
}
