import React from 'react';
import { createRoot } from 'react-dom/client';

import App from './App';
import './index.css';
import { ChakraProvider } from '@chakra-ui/react';
import '@fontsource/manrope/700.css';
import '@fontsource/manrope/400.css';
import theme from '../../util/theme';

const container = document.getElementById('app-container');
const root = createRoot(container); // createRoot(container!) if you use TypeScript

root.render(
  <ChakraProvider theme={theme}>
    <App />
  </ChakraProvider>
);
