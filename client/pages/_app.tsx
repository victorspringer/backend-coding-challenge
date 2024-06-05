import * as React from 'react';
import Head from 'next/head';
import { AppProps } from 'next/app';
import { AppCacheProvider } from '@mui/material-nextjs/v14-pagesRouter';
import { ThemeProvider } from '@mui/material/styles';
import CssBaseline from '@mui/material/CssBaseline';
import theme from '../src/theme';
import Navigation from '../src/components/Navigation';
import { Container } from '@mui/material';
import Box from '@mui/material/Box';
import '../public/styles/theme.css';
import Footer from '../src/components/Footer';

const showNavigation = (pathname: string): boolean => {
  switch (pathname) {
    case "/":
    case "/profile/[username]":
    case "/movies":
    case "/users":
      return true;
  }
  return false;
};

export default function MyApp(props: AppProps) {
  const { Component, pageProps } = props;
  
  return (
    <AppCacheProvider {...props}>
      <Head>
        <meta name="viewport" content="initial-scale=1, width=device-width" />
      </Head>
      <ThemeProvider theme={theme}>
        {/* CssBaseline kickstart an elegant, consistent, and simple baseline to build upon. */}
        <CssBaseline />
        {showNavigation(props.router.pathname) && <Navigation />}
        <Container maxWidth="xl">
          <Box
            sx={{
              my: 4,
              display: 'flex',
              flexDirection: 'column',
              justifyContent: 'center',
              alignItems: 'center',
            }}
          >
            <Component {...pageProps} />
          </Box>
        </Container>
        <Footer />
      </ThemeProvider>
    </AppCacheProvider>
  );
}
