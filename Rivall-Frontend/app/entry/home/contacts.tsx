import { Image } from '@/components/ui/image';
import { Box } from '@/components/ui/box';
import { Text } from '@/components/ui/text';
import { HStack } from '@/components/ui/hstack';
import { VStack } from '@/components/ui/vstack';
import { Icon } from '@/components/ui/icon';
import { Badge, BadgeText } from '@/components/ui/badge';
import { Avatar, AvatarFallbackText, AvatarImage } from '@/components/ui/avatar';
import { BackgroundGradientWrapper } from '@/components/BackgroundGradientWrapper';
import { FlatList, View, Pressable } from 'react-native';
import { useRouter } from 'expo-router';
import { Card } from '@/components/ui/card';
import data from '@/api/contacts.json';
import { RivallBlackCrown } from '@/components/RivallBlackCrown';
import { Crown } from 'lucide-react-native';

export default function ContactsScreen() {

  const contacts = data;

  return (
    <BackgroundGradientWrapper>
      <FlatList
        data={contacts}
        renderItem={({ item }) => (
          <ContactCard contact={item} />
        )}
        keyExtractor={item => item._id}
      >
      </FlatList>
    </BackgroundGradientWrapper>
  )
}

interface Contact {
  _id: string;
  group_name: string;
  num_members: number;
  is_challenge: boolean;
  challenge_start_date: string;
  unread_count: number;
  last_message: {
    _id: string;
    first_name: string;
    last_name: string;
    avatar_image: string;
    avatar_color: string;
    message: string;
    timestamp: string;
    unread: boolean;
  }
}

function ContactCard({ contact }: { contact: Contact }) {

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

  const router = useRouter();

  return (
    <Pressable onPress={() => router.push(`entry/chat/${contact._id}`)}>
      <Card className="my-2 mx-4">
        <VStack>
          {contact.is_challenge ? (
            <HStack className="gap-2">
              <Icon as={Crown} size="lg" className="my-auto text-black"/>
              <Badge size="lg" className="bg-secondary-500">
                <BadgeText className="font-extrabold">GROUP</BadgeText>
              </Badge>
            </HStack>) : null
          }
          <HStack>
            <Avatar className={`w-14 h-14 mr-4 my-auto`} style={{ backgroundColor: contact.last_message.avatar_color }}>
              <AvatarFallbackText className="text-white">
                {contact.last_message.first_name[0] + ' ' + contact.last_message.last_name[0]}
              </AvatarFallbackText>
              <AvatarImage
                source={{ uri: contact.last_message.avatar_image }}
                className="w-14 h-14 rounded-full"
              />
            </Avatar>
            <VStack className="flex-1 justify-start">
              <HStack className="gap-2 mb-2">
                <Text className="font-bold text-black mt-auto">{contact.last_message.first_name} {contact.last_message.last_name}</Text>
                <View className="w-fit ml-auto">
                  {contact.unread_count > 0 ? (<Badge className="bg-sky-500 rounded-full text-center">
                    <BadgeText className="text-white font-bold text-sm mx-auto">{ contact.unread_count }</BadgeText>
                  </Badge>) : null}
                </View>
              </HStack>
              <Text className="text-wrap">{ contact.last_message.message.slice(0, 70) + (contact.last_message.message.length > 70 ? '...' : '') }</Text>
              <Text className="ml-auto">{ formatTimestamp(contact.last_message.timestamp) }</Text>
            </VStack>
          </HStack>
        </VStack>
      </Card>
    </Pressable>
  )
}