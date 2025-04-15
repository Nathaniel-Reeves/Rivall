import { Box } from '@/components/ui/box';
import { Text } from '@/components/ui/text';
import { Card } from '@/components/ui/card';
import { HStack } from '@/components/ui/hstack';
import { Avatar, AvatarFallbackText, AvatarImage } from '@/components/ui/avatar';
import { useUserStore } from '@/global-store/user_store';
import { Message, User, ChatMembers } from '@/types';

export default function MessageBox({ message, chatMembers }: { message: Message, chatMembers: ChatMembers }) {

  const user_store = useUserStore((state: any) => state.user);
  let isUserMessage = false;
  if (user_store._id === message.user_id) {
    isUserMessage = true;
  }

  const user = (chatMembers[message.user_id]);

  return (
    <>
      {
        isUserMessage ? <UserMessage message={message} user={user} /> : <ContactMessage message={message} user={user} />
      }
    </>
  )
}

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

function ContactMessage({ message, user }: { message: Message, user: User }) {
  return (
    <HStack className="m-4">
      <Box className="w-14 h-14 mr-4 my-auto shadow-md shadow-black rounded-full">
        <Avatar className={`w-14 h-14 mr-4 my-auto shadow-md shadow-black`} style={{ backgroundColor: "#555555" }}>
          <AvatarFallbackText className="text-white">
            { user.first_name + ' ' + user.last_name }
          </AvatarFallbackText>
        </Avatar>
      </Box>
      <Box>
        <Text className="typography-900">{ user.first_name + ' ' + user.last_name }</Text>
        <Card className="max-w-[80%] w-fit p-3 shadow-md shadow-black">
          <Text>{message.message_data}</Text>
        </Card>
        <Text className="typography-900 mt-2">{formatTimestamp(message.timestamp)}</Text>
      </Box>
    </HStack>
  )
}

function UserMessage({ message, user }: { message: Message, user: any }) {
  return (
    <HStack className="justify-end m-4">
      <Box>
        <Text className="typography-900 text-right ">{ user.first_name + ' ' + user.last_name }</Text>
        <Card className="max-w-[80%] w-fit p-3 shadow-md shadow-black ml-auto bg-info-800">
          <Text className="text-right text-white">{message.message_data}</Text>
        </Card>
        <Text className="typography-900 mt-2 text-right">{formatTimestamp(message.timestamp)}</Text>
      </Box>
      <Box className="w-14 h-14 ml-4 my-auto shadow-md shadow-black rounded-full">
        <Avatar className={`w-14 h-14 mr-4 my-auto shadow-md shadow-black`} style={{ backgroundColor: "#555555" }}>
          <AvatarFallbackText className="text-white">
            { user.first_name + ' ' + user.last_name }
          </AvatarFallbackText>
        </Avatar>
      </Box>
    </HStack>
  )
}

