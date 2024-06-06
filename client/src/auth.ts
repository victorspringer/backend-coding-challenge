const RefreshToken = async (refreshToken: string) => {
    const response = await fetch(`${process.env.NEXT_PUBLIC_WEBAPP_URL}/api/refreshToken`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({refreshToken}),
    });
    const data = await response.json()
    return data;
};

export default RefreshToken;
