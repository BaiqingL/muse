import React, { useState } from 'react';
import {
  Box,
  Button,
  ChakraProvider,
  Heading,
  Text,
  Textarea,
} from '@chakra-ui/react';
import './Popup.css';

const SAMPLE_APP_IDEAS = [
  'An app that helps you discover random inspirational quotes',
  'An app that allows you to efficiently manage and track tasks and to-do lists',
  'An app that enables you to find recipes by searching available ingredients',
  'An app that provides personalized book recommendations',
  'An app that tracks and monitors fitness activities and progress',
  'An app that helps you master a new language through interactive flashcards',
  'An app that effortlessly tracks and manages daily expenses',
  'An app that lets you create countdowns for exciting upcoming events',
  'An app that helps you explore movies and share your reviews',
  'An app that helps you relax and focus with a guided meditation timer',
];

const SELECTED_SAMPLE =
  SAMPLE_APP_IDEAS[Math.floor(Math.random() * SAMPLE_APP_IDEAS.length)];

const Popup = () => {
  const [prompt, setPrompt] = useState();

  return (
    <Box>
      <Heading size="md">Create your new dream app</Heading>
      <Textarea
        placeholder={SELECTED_SAMPLE}
        value={prompt}
        onChange={(e) => setPrompt(e.target.value)}
      />
      <Button>Create</Button>
    </Box>
  );
};

export default Popup;
