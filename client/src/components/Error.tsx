import * as React from 'react';
import theme from '../theme';
import { Box, Button, Typography } from '@mui/material';
import NotFoundError from '../../public/svg/notFoundError.svg';
import InternalServerError from '../../public/svg/internalServerError.svg';
import { useRouter } from 'next/router';

type Props = {
  code: number;
};

export default function Error({ code }: Props) {
    const router = useRouter();

    switch(code) {
        case 400:
          break;
        case 401:
          break;
        case 404:
          return (
            <Box textAlign='center'>
              <NotFoundError />
              <Typography mt={6} mb={8} variant="h4" color={theme.palette.grey[500]}>
                The page you are trying to access was not found.
              </Typography>
              <Button onClick={() => router.back()} type='button' variant='contained' size='large' sx={{ display: 'block', margin: '32px auto' }}>
                Click here to go back
              </Button>
            </Box>
          );
        default:
          return (
            <Box textAlign='center'>
              <InternalServerError />
              <Typography mb={8} variant="h4" color={theme.palette.grey[500]}>
                Oops... something went wrong. Sorry for the inconvenience.
              </Typography>
              <Button onClick={() => router.back()} type='button' variant='contained' size='large' sx={{ display: 'block', margin: '32px auto' }}>
                Click here to go back
              </Button>
            </Box>
          );
    }
}
