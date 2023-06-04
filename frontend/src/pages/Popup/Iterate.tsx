import {
  Heading,
  Stack,
  Text,
  Textarea,
  Button,
  Image,
  Box,
} from '@chakra-ui/react';
import React from 'react';

// @ts-ignore
import logo from '../../assets/img/logo.svg';

const Iterate = () => {
  return (
    <Stack spacing={4} p={4}>
      <Box alignSelf="center">
        <Image src={logo} alt="logo" maxH="50px" />
      </Box>
      <Text color="gray">
        Provide feedback on the generated prototype until you're satisfied
      </Text>
      <Textarea placeholder="What would you like to change?" />
      <Button colorScheme="purple">Submit</Button>
    </Stack>
  );
};

export default Iterate;
