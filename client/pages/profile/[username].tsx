import * as React from 'react';
import Typography from '@mui/material/Typography';
import { Avatar, Box, Card, CardContent, CardMedia, Rating } from '@mui/material';
import fetch from 'isomorphic-fetch';
import { GetServerSideProps } from 'next';
import theme from '../../src/theme';
import StarIcon from '@mui/icons-material/Star';
import Error from '../../src/components/Error';
import CircularProgress from '@mui/material/CircularProgress';
import IsAuthenticated from '../../src/auth';
import { GetCookie } from '../../src/cookies';

type User = {
    id: string;
    name: string;
    username: string;
    md5Password: string;
    picture: string;
};

type Rating = {
    user: User;
    movie: Movie;
    value: number;
};

type Movie = {
    id: string;
    title: string;
    poster: string;
};

type Props = {
    user?: User;
    ratings?: Rating[];
    error?: Error;
};

type Error = {
    code: number;
};

export const getServerSideProps: GetServerSideProps = async ({ req, res, query }) => {
    const isAuthenticated = await IsAuthenticated(req, res);
    if (!isAuthenticated) return {
        redirect: {
            destination: '/signin',
            permanent: false,
        },
    };

    const { username } = query;
    const accessToken = GetCookie(req, "MRSAccessToken");
    const props: Props = {};


    const userResponse = await fetch(`${process.env.NEXT_PUBLIC_WEBAPP_URL}/api/user/${username}?accessToken=${accessToken}`);
    const userData = await userResponse.json();

    if (userData.error) {
        console.log(userData.error);
        props.error = { code: userData.statusCode }
        return { props }
    }

    props.user = userData.response;

    const ratingsResponse = await fetch(`${process.env.NEXT_PUBLIC_WEBAPP_URL}/api/rating/${username}?accessToken=${accessToken}`);

    const ratingsData = await ratingsResponse.json()

    if (ratingsData.error) {
        console.log(ratingsData.error);
        props.error = { code: ratingsData.statusCode }
        return { props }
    }

    const ratings = await Promise.all(
        ratingsData.response.map(async (rating: any) => {
            const movieResponse = await fetch(`${process.env.NEXT_PUBLIC_WEBAPP_URL}/api/movie/${rating.movieId}?accessToken=${accessToken}`);

            const movieData = await movieResponse.json();

            if (movieData.error) {
                console.log(movieData.error);
                return null;
            }

            return {
                user: userData.response,
                movie: movieData.response,
                value: rating.value,
            };
        }).filter((rating: Rating) => rating !== null)
    );

    props.ratings = ratings;

    return { props };
};

const updateRating = async (value: number, rating?: Rating) => {
    if (!rating) return;

    const response = await fetch(`/api/updateRating?accessToken=${localStorage.getItem("accessToken")}`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({
            userId: rating.user.id,
            movieId: rating.movie.id,
            value,
        }),
    });

    if (response.ok) {
        return null;
    } else {
        const data = await response.json();
        return data.statusCode;
    }
};

const labels: { [index: string]: string } = {
    0.5: 'Terrible',
    1: 'Very Poor',
    1.5: 'Poor',
    2: 'Below Average',
    2.5: 'Average',
    3: 'Above Average',
    3.5: 'Good',
    4: 'Very Good',
    4.5: 'Great',
    5: 'Masterpiece',
};

export default function Profile({ user, ratings, error }: Props) {
    if (error) {
        return <Error code={error.code} />;
    }

    if (!user) {
        return (
            <Box my={20}>
                <CircularProgress color='primary' />
            </Box>
        );
    }

    const [loggedInUser, setLoggedInUser] = React.useState<string|null>("");
    React.useEffect(() => {
      if (typeof window !== "undefined") {
        setLoggedInUser(localStorage.getItem("loggedInUser"));
      }
    }, []);

    const firstName = user.name.split(" ")[0];

    const [values, setValues] = React.useState<number[]>(ratings?.map(rating => rating.value) || []);
    const [hover, setHover] = React.useState<number[]>(ratings?.map(rating => rating.value) || []);

    const handleChange = (index: number) => async (event: React.ChangeEvent<{}>, newValue: number | null) => {
        if (newValue !== null) {
            const newValues = [...values];
            newValues[index] = newValue;
            setValues(newValues);
            error = await updateRating(newValue, ratings ? ratings[index] : undefined);
        }
    };

    const handleHover = (index: number) => (event: React.ChangeEvent<{}>, newHover: number | null) => {
        if (newHover !== null) {
            const newHovers = [...values];
            newHovers[index] = newHover;
            setHover(newHovers);
        }
    };

    return (
        <Box
            sx={{
                display: 'flex',
                flexDirection: 'row',
                flexFlow: 'row wrap',
                justifyContent: 'center',
                gap: '32px',
                width: '100%',
            }}
        >
            <Card className="card" variant='outlined' sx={{ width: 275 }}>
                <CardMedia
                    component="div"
                    sx={{
                        height: 140,
                        backgroundImage: "linear-gradient(to top right,#d10000,#f60439,#b105f4)",
                    }}
                />
                <CardContent>
                    <Avatar
                        className='profile-avatar profile-avatar-l'
                        alt={user.name}
                        src={user.picture}
                    />
                    <Typography gutterBottom variant="h4" component="div">
                        {user.name}
                    </Typography>
                    <Typography gutterBottom variant="h6" component="div" color={theme.palette.grey[500]}>
                        @{user.username}
                    </Typography>
                </CardContent>
            </Card>
            <Card className="card" sx={{ flexGrow: 1, minWidth: 320, maxWidth: 'calc(100% - 275px - 32px)' }}>
                <CardContent>
                    <Typography mb={2} variant="h5" component="div" align='center'>
                        {firstName}{firstName.endsWith('s') ? "'" : "'s"} Ratings
                    </Typography>
                    <Box
                        sx={{
                            display: 'flex',
                            flexDirection: 'row',
                            flexFlow: 'row wrap',
                            justifyContent: 'center',
                            gap: '16px',
                        }}>
                        {
                            ratings?.map((rating, i) => {
                                return (
                                    <Card key={i} className='movie-card' variant='outlined'>
                                        <CardMedia
                                            component="div"
                                            sx={{
                                                height: 375,
                                                backgroundImage: `url(${rating.movie.poster})`,
                                            }}
                                        />
                                        <CardContent>
                                            <Typography gutterBottom variant="h6" fontSize={18}>
                                                {rating.movie.title}
                                            </Typography>
                                            <Rating
                                                readOnly={user.username !== loggedInUser}
                                                precision={0.5}
                                                value={values[i]}
                                                onChange={handleChange(i)}
                                                onChangeActive={handleHover(i)}
                                                emptyIcon={<StarIcon style={{ opacity: 0.55 }} fontSize="inherit" />}
                                            />
                                            <Typography mt={1} variant="body1" fontSize={16} fontStyle="italic" align='right'>
                                                "{labels[hover[i] !== -1 ? hover[i] : values[i]]}"
                                            </Typography>
                                        </CardContent>
                                    </Card>
                                );
                            })
                        }
                    </Box>
                </CardContent>
            </Card>
        </Box>
    );
}
