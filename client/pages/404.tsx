import * as React from 'react';
import Typography from '@mui/material/Typography';
import Error from '../src/components/Error';

export default function Custom404() {
  return <Error code={404} />;
}
