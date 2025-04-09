import { BackgroundGradientWrapper } from '@/components/BackgroundGradientWrapper';
import { FlatList, View } from 'react-native';
import { useEffect } from 'react';
import { Stack, useLocalSearchParams, useNavigation } from "expo-router";
import { Input, InputField, InputSlot, InputIcon } from '@/components/ui/input';
import { Textarea, TextareaInput } from '@/components/ui/textarea';
import { Box } from '@/components/ui/box';
import { ArrowUp } from 'lucide-react-native';
import { Icon } from '@/components/ui/icon';
import { Text } from '@/components/ui/text';
import { HStack } from '@/components/ui/hstack';
import { Button, ButtonIcon } from '@/components/ui/button';

export default function ChatScreen() {
  const { id } = useLocalSearchParams<{ id: string }>();
  console.log(id);
  
  const messageData = {
    group_members: {
      "someID": {
        first_name: "John",
        last_name: "Doe",
        avatar_image: "https://example.com/image.jpg",
        avatar_color: "#000000"
      },
      "anotherID": {
        first_name: "Jane",
        last_name: "Doe",
        avatar_image: "https://example.com/image.jpg",
        avatar_color: "#000000"
      }
    },
    messages: [
      {
        _id: "someID",
        message: "Hello, world!",
        timestamp: "2022-01-01T00:00:00.000Z"
      },
      {
        _id: "anotherID",
        message: "Hi, there!",
        timestamp: "2022-01-01T00:00:01.000Z"
      }
    ]
  }

  const nav = useNavigation();

  useEffect(() => {
    nav.setOptions({
      title: 'Chat'
    });
  }, []);

  return (
    <BackgroundGradientWrapper>
      <Stack.Screen
        options={{
          title: messageData.group_members["someID"].first_name + " " + messageData.group_members["someID"].last_name,
          headerShown: true
        }}
      />
      <FlatList
        data={messageData.messages}
        renderItem={({ item }) => (
          <Message message={item} />
        )}
        keyExtractor={item => item._id}
      >
      </FlatList>
      <Box className="bg-neutral-300 w-full h-1"></Box>
      <Box className="bottom-0 w-full h-16 bg-white p-4">
        <HStack className="flex-1 justify-end gap-2">
          <Textarea className="rounded-2xl w-3/4 h-10" size="md">
            <TextareaInput className="align-top" placeholder="Type a message..."></TextareaInput>
          </Textarea>
          <Button className="rounded-full bg-info-800 w-10 h-10">
            <ButtonIcon as={ArrowUp} size="2xl" className="text-white"></ButtonIcon>
          </Button>
        </HStack>
      </Box>
    </BackgroundGradientWrapper>
  )
}

interface Message {
  _id: string;
  message: string;
  timestamp: string;
}

function Message({ message }: { message: Message }) {

  function formatTimestamp(timestamp: string): string {
    const now = new Date();
    const time = new Date(timestamp);
    const diffInSeconds = Math.floor((now.getTime() - time.getTime()) / 1000);

    if (diffInSeconds < 60) {
      return `${diffInSeconds}s`;
    } else if (diffInSeconds < 3600) {
      const minutes = Math.floor(diffInSeconds / 60);
      return `${minutes}m`;
    } else if (diffInSeconds < 86400) {
      return time.toLocaleTimeString([], { hour: 'numeric', minute: '2-digit', hour12: true });
    } else {
      return time.toLocaleDateString();
    }
  }

  return (
    <View>
      <Text>{message.message}</Text>
      <Text>{formatTimestamp(message.timestamp)}</Text>
    </View>
  )
}