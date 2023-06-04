import React, { useState } from 'react';
import {
  Box,
  Button,
  Heading,
  Image,
  Progress,
  Select,
  Stack,
  Text,
  Textarea,
} from '@chakra-ui/react';
import './Home.css';
import { api } from '../../util/api';

// @ts-ignore
import logo from '../../assets/img/logo.svg';

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

const APP_URL = 'http://localhost:5173';

const Popup = () => {
  const [prompt, setPrompt] = useState<string>();
  const [framework, setFramework] = useState<(typeof UI_FRAMEWORKS)[0]>();
  const [loading, setLoading] = useState<boolean>(false);
  const [envContents, setEnvContents] = useState<string>('');

  const handleCreate = async () => {
    if (!prompt || !framework) {
      return;
    }

    try {
      setLoading(true);

      const apiKey = await chrome.storage.local.get(['apiKey']);

      if (!apiKey) {
        return;
      }

      // const resColdStart = await api.post('/api/coldStart', {
      //   framework: framework.id,
      //   useCase: prompt,
      //   apiKey: apiKey.apiKey,
      // });

      const resCheckEnv = await api.get('/api/getFile', {
        params: {
          filename: '.env',
        },
      });

      console.log(resCheckEnv);

      if (resCheckEnv.data.exist) {
        const contents = resCheckEnv.data.content;
        console.log(contents);

        setEnvContents(contents);
      }

      // await chrome.tabs.create({
      //   url: APP_URL,
      // });
    } catch (err) {
      console.log(err);
    }

    setLoading(false);
  };

  return (
    <Stack spacing={4} p={4}>
      <Box alignSelf="center">
        <Image src={logo} alt="logo" maxH="50px" />
      </Box>
      <Text color="gray">
        Start with an idea, and we'll generate a prototype
      </Text>
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
      <Button
        onClick={handleCreate}
        colorScheme="purple"
        alignSelf="end"
        isLoading={loading}
      >
        Create
      </Button>

      {loading && <Progress size="xs" isIndeterminate />}

      {envContents && (
        <Stack>
          <Heading size="md">.env</Heading>
          <Text>{envContents}</Text>
        </Stack>
      )}
    </Stack>
  );
};

export default Popup;
