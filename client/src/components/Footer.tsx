import * as React from 'react';
import theme from '../theme';
import { Link, Typography } from '@mui/material';
import FavoriteIcon from '@mui/icons-material/Favorite';

export default function Footer() {
  return (
    <footer>
      <Typography align='center' variant="subtitle1" fontSize={13} color={theme.palette.grey[500]}>
        Made with <FavoriteIcon
          color='primary'
          sx={{ fontSize: 'smaller' }}
        /> for <Link href='https://www.thermondo.de/' target='_blank' color={theme.palette.primary.light} underline='none'>thermondo</Link>
      </Typography>
      <Typography mb={2} align='center' variant="subtitle1" fontSize={13} color={theme.palette.grey[500]}>
        Author: <Link href='https://github.com/victorspringer/backend-coding-challenge' target='_blank' color={theme.palette.secondary.light} underline='none'>Victor Springer</Link>
      </Typography>
    </footer>
  );
}
