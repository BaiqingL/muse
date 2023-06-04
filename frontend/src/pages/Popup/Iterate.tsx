import { Heading, Stack, Text, Textarea, Button } from '@chakra-ui/react';
import React from 'react';

const Iterate = () => {
  return (
    <Stack spacing={4} p={4}>
      <Heading size="lg" alignSelf="center">
        Iterate
      </Heading>
      <Text color="gray">
        Provide feedback on the generated prototype until you're satisfied
      </Text>
      <Textarea placeholder="What would you like to change?" />
      <Button colorScheme="purple">Submit</Button>
    </Stack>
  );
};

export default Iterate;
