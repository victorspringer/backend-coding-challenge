import * as React from 'react';
import Typography from '@mui/material/Typography';
import { Alert, Box, Button, Card, CardContent, CardMedia, Checkbox, FormControlLabel, FormGroup, Link, TextField } from '@mui/material';
import theme from '../src/theme';
import { md5 } from 'js-md5';
import { GetServerSideProps } from 'next';
import fetch from 'isomorphic-fetch';
import cookie from 'cookie';
import { useRouter } from 'next/router';

export const getServerSideProps: GetServerSideProps = async ({ req, res, query }) => {
    const cookies = req.headers.cookie ? cookie.parse(req.headers.cookie) : null;
    const isAuthenticated = cookies && cookies["MRSRefreshToken"];
    
    if (isAuthenticated) {
        return {
            redirect: {
                destination: `/profile/victorspringer`,
                permanent: false,
            },
        };
    }

    return { props: {} };
}

export default function SignIn() {
    const [username, setUsername] = React.useState("");
    const [password, setPassword] = React.useState("");
    const [rememberMe, setRememberMe] = React.useState(true);
    const [error, setError] = React.useState(false);
    const router = useRouter();

    const handleSubmit = async (e: React.MouseEvent<HTMLAnchorElement>) => {
        e.preventDefault();
        setError(false);

        const response = await fetch(`api/signin`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                flow: rememberMe ? 'rememberMe' : 'websiteSession',
                md5Password: md5(password),
                username: username,
            }),
        });
        const data = await response.json()
        if (data.statusCode == 200) {
            setError(false);
            router.push(`/profile/${username}`)
        } else {
            setError(true);
        }
    };

    return (
        <Box>
            <Alert sx={{ visibility: !error ? 'hidden' : 'visible', marginTop: '56px' }} severity="error">Wrong credentials.</Alert>
            <Card sx={{ width: '20vw', minWidth: '280px', margin: '8px auto' }}>
                <CardMedia
                    component="div"
                    sx={{
                        height: 100,
                        backgroundImage: "linear-gradient(to top right,#d10000,#f60439,#b105f4)",
                    }}
                >
                    <Typography paddingTop={4} variant="h5" fontWeight='bold' component="h1" color='#FFFFFF' textAlign='center'>
                        Sign in to your account
                    </Typography>
                </CardMedia>
                <CardContent>
                    <FormGroup>
                        <TextField
                            label="Username"
                            required
                            variant="outlined"
                            fullWidth
                            margin='normal'
                            value={username}
                            onChange={(e) => setUsername(e.target.value)}
                        />
                        <TextField
                            label="Password"
                            required
                            variant="outlined"
                            type='password'
                            fullWidth
                            margin='normal'
                            value={password}
                            onChange={(e) => setPassword(e.target.value)}
                        />
                        <FormControlLabel control={<Checkbox checked={rememberMe} onClick={() => setRememberMe(!rememberMe)} />} label="Remember me" />
                        <Button type='submit' href="" variant='contained' size='large' onClick={handleSubmit}>Sign in</Button>
                    </FormGroup>
                    <Typography variant='subtitle2' fontWeight={400} color={theme.palette.grey[700]} marginTop={2} align='right'>
                        Don't have an account yet? <Link href='/signup'>Sign up</Link>
                    </Typography>
                </CardContent>
            </Card>
        </Box>
    );
}
