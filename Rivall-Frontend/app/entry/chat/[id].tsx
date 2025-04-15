import { BackgroundGradientWrapper } from '@/components/BackgroundGradientWrapper';
import { FlatList, View } from 'react-native';
import { useEffect, useState } from 'react';
import { Stack, useLocalSearchParams, useNavigation } from "expo-router";
import { Input, InputField, InputSlot, InputIcon } from '@/components/ui/input';
import { Textarea, TextareaInput } from '@/components/ui/textarea';
import { Box } from '@/components/ui/box';
import { ArrowUp } from 'lucide-react-native';
import { Icon } from '@/components/ui/icon';
import { Text } from '@/components/ui/text';
import { HStack } from '@/components/ui/hstack';
import { Button, ButtonIcon } from '@/components/ui/button';
import MessageBox from '@/components/MessageBox';
import uuid from 'react-native-uuid';
import { useUserStore } from '@/global-store/user_store';
import { useQuery } from '@tanstack/react-query';
import { getChat } from '@/api/contact';
import { useWebSockets } from '@/hooks/useWebSocket';

export default function ChatScreen() {
  const { id } = useLocalSearchParams<{ id: string }>();
  console.log(id);

  const user = useUserStore((state: any) => state.user);
  const [otherUser, setOtherUser] = useState<any>({});
  const access_token = useUserStore((state: any) => state.access_token);
  const [ messageContent, setMessageContent ] = useState<string>("");

  const handleReceivedMessage = (message: any) => {
    console.log("Received message: ", message);
    setMessageData((prevState: any) => ({
      ...prevState,
      messages: [...prevState.messages, message]
    }));
    setMessageContent("");
  }

  const { sendMessage } = useWebSockets({receivedMessage: handleReceivedMessage});
  
  const [messageData, setMessageData] = useState<any>({
    group_members: {},
    messages: []
  });

  // Get User Data using auth token
  const { data, isLoading, error } = useQuery({
    queryKey: ['getChat', id],
    queryFn: () => getChat(user._id, access_token, id),
    retryDelay: attempt => Math.min(attempt > 1 ? 2 ** attempt * 1000 : 1000, 30 * 1000),
  });

  useEffect(() => {
    console.log("Chat Data: ", data);
    if (data) {
      console.log(JSON.stringify(data, null, 2));
      setMessageData(data.data);
      const otherUserID = Object.keys(data.data.group_members).find(key => key !== user._id);
      setOtherUser(data.data.group_members[otherUserID]);
    }
  }, [data]);

  if (isLoading) {
    return (
      <View className="flex-1 justify-center w-80 mx-auto">
        <Text className="text-typography-800 text-2xl font-medium text-pretty text-center mb-20">Loading...</Text>
      </View>
    )
  }

  if (error) {
    console.error(error)
    return (
      <View className="flex-1 justify-center w-80 mx-auto">
        <Text className="text-typography-800 text-2xl font-medium text-pretty text-center mb-20">Error loading chat</Text>
      </View>
    )
  }

  const handleSendMessage = () => {
    if (messageContent.trim() === "") {
      return;
    }
    const newMessage = {
      _id: uuid.v4(),
      user_id: user._id,
      receiver_id: otherUser._id,
      message_data: messageContent,
      timestamp: new Date().toISOString(),
      message_type: "text"
    }
    sendMessage(newMessage, id);
    setMessageData((prevState: any) => ({
      ...prevState,
      messages: [...prevState.messages, newMessage]
    }));
    setMessageContent("");
  }

  return (
    <BackgroundGradientWrapper>
      <Stack.Screen
        options={{
          headerShown: true,
          headerTitle: otherUser.first_name + " " + otherUser.last_name,
        }}
      />
      <Box className="bg-neutral-300 w-full h-1">
        <Text className="text-center text-neutral-500 text-xs mt-2">Today</Text>
      </Box>
      <FlatList
        data={messageData?.messages}
        renderItem={({ item }) => (
          <MessageBox message={item} chatMembers={messageData.group_members} />
        )}
        keyExtractor={item => item._id}
      >
      </FlatList>
      <Box className="bg-neutral-300 w-full h-1"></Box>
      <Box className="bottom-0 w-full h-16 bg-white p-4">
        <HStack className="flex-1 justify-end gap-2">
          <Textarea className="rounded-2xl w-3/4 h-10" size="md">
            <TextareaInput className="align-top" placeholder="Type a message..." onChangeText={(text) => setMessageContent(text)}></TextareaInput>
          </Textarea>
          <Button className="rounded-full bg-info-800 w-10 h-10" onPress={handleSendMessage}>
            <ButtonIcon as={ArrowUp} size="2xl" className="text-white"></ButtonIcon>
          </Button>
        </HStack>
      </Box>
    </BackgroundGradientWrapper>
  )
}