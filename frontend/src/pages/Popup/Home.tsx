import React, { useEffect, useState } from 'react';
import {
  Box,
  Button,
  ChakraProvider,
  Heading,
  Select,
  Stack,
  Text,
  Textarea,
} from '@chakra-ui/react';
import './Home.css';
import { api } from '../../util/api';

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

const UI_FRAMEWORKS = [
  {
    name: 'Chakra UI',
    id: 'chakra-ui',
  },
  {
    name: 'Material UI',
    id: 'material-ui',
  },
  {
    name: 'Tailwind CSS',
    id: 'tailwind-css',
  },
  {
    name: 'Bootstrap',
    id: 'bootstrap',
  },
  {
    name: 'Ant Design',
    id: 'ant-design',
  },
];

const SELECTED_SAMPLE =
  SAMPLE_APP_IDEAS[Math.floor(Math.random() * SAMPLE_APP_IDEAS.length)];

const APP_URL = 'https://google.com/';

const Popup = () => {
  const [prompt, setPrompt] = useState<string>();
  const [framework, setFramework] = useState<(typeof UI_FRAMEWORKS)[0]>();

  const handleCreate = async () => {
    const res = await api.post('/api/userPrompt', { prompt });

    console.log(res.data);

    await chrome.tabs.create({
      url: APP_URL,
    });
  };

  return (
    <Stack spacing={4} p={4}>
      <Heading size="md" alignSelf="center">
        Create your new dream app
      </Heading>
      <Textarea
        placeholder={SELECTED_SAMPLE}
        value={prompt}
        onChange={(e) => setPrompt(e.target.value)}
      />
      <Select
        placeholder="Select UI framework"
        onChange={(e) =>
          setFramework(UI_FRAMEWORKS.find((f) => f.id === e.target.value))
        }
      >
        {UI_FRAMEWORKS.map((framework) => (
          <option key={framework.id} value={framework.id}>
            {framework.name}
          </option>
        ))}
      </Select>
      <Button onClick={handleCreate}>Create</Button>
      <Text>
        {prompt} {framework?.name}
      </Text>
    </Stack>
  );
};

export default Popup;
