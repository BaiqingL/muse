import React, { useState } from 'react';
import './Options.css';
import {
  Box,
  Button,
  Flex,
  FormControl,
  FormLabel,
  Heading,
  Input,
  Stack,
  useToast,
} from '@chakra-ui/react';

const Options = () => {
  const toast = useToast();
  const [apiKey, setApiKey] = useState<string>('');

  const handleSave = async () => {
    toast({
      title: 'Saved!',
      description: 'Your API key has been saved.',
      status: 'success',
    });

    await chrome.storage.local.set({
      apiKey,
    });

    await chrome.storage.local.get(['apiKey'], (result) => {
      console.log('Value currently is ' + result.apiKey);
    });
  };

  return (
    <Flex justify="center" w="full" h="100vh">
      <Stack alignSelf="center" maxW="1000px" w="full">
        <FormControl>
          <FormLabel>OpenAI API Key</FormLabel>
          <Input
            type="password"
            placeholder="sk-your-api-key"
            onChange={(e) => setApiKey(e.target.value)}
          />
        </FormControl>
        <Button colorScheme="purple" alignSelf="center" onClick={handleSave}>
          Save
        </Button>
      </Stack>
    </Flex>
  );
};

export default Options;
