import {
  Heading,
  Stack,
  Text,
  Textarea,
  Button,
  Image,
  Box,
} from '@chakra-ui/react';
import React, { useState } from 'react';

// @ts-ignore
import logo from '../../assets/img/logo.svg';
import { api } from '../../util/api';

const Iterate = () => {
  const [prompt, setPrompt] = useState<string>();
  const [loading, setLoading] = useState<boolean>(false);

  const handleSubmit = async () => {
    try {
      setLoading(true);
      const res = await api.post('/api/iterate', {
        prompt,
      });
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
        Provide feedback on the generated prototype until you're satisfied
      </Text>
      <Textarea
        placeholder="What would you like to change?"
        value={prompt}
        onChange={(e) => setPrompt(e.target.value)}
      />
      <Button colorScheme="purple" onClick={handleSubmit} isLoading={loading}>
        Submit
      </Button>
    </Stack>
  );
};

export default Iterate;
