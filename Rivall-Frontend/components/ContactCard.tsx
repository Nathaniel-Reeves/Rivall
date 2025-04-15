import { Text } from '@/components/ui/text';
import { HStack } from '@/components/ui/hstack';
import { VStack } from '@/components/ui/vstack';
import { Icon } from '@/components/ui/icon';
import { Badge, BadgeText } from '@/components/ui/badge';
import { Avatar, AvatarFallbackText } from '@/components/ui/avatar';
import { View, Pressable } from 'react-native';
import { useRouter } from 'expo-router';
import { Card } from '@/components/ui/card';
import { Crown } from 'lucide-react-native';
import { Contact } from '@/types';

export default function ContactCard({ contact }: { contact: Contact }) {

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
  
    function avatarColor() {
      const colors = ['#FF5733', '#33FF57', '#3357FF', '#FF33A1', '#A133FF'];
      const index = Math.floor(Math.random() * colors.length);
      return colors[index];
    }
  
    return (
      <Pressable onPress={() => router.push(`entry/chat/${contact.direct_message._id}`)}>
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
              <Avatar className={`w-14 h-14 mr-4 my-auto`} style={{ backgroundColor: avatarColor() }}>
                <AvatarFallbackText className="text-white">
                  {contact.first_name[0] + ' ' + contact.last_name[0]}
                </AvatarFallbackText>
                {/* <AvatarImage
                  source={{ uri: contact.avatar_image }}
                  className="w-14 h-14 rounded-full"
                /> */}
              </Avatar>
              <VStack className="flex-1 justify-start">
                <HStack className="gap-2 mb-2">
                  <Text className="font-bold text-black mt-auto">{contact.first_name} {contact.last_name}</Text>
                  <View className="w-fit ml-auto">
                    {contact.unread_count > 0 ? (<Badge className="bg-sky-500 rounded-full text-center">
                      <BadgeText className="text-white font-bold text-sm mx-auto">{ contact.unread_count }</BadgeText>
                    </Badge>) : null}
                  </View>
                </HStack>
                  { contact.direct_message.messages.length === 0 ?
                    (<Text className="text-wrap">No Messages Yet</Text>) : 
                    <>
                      <Text className="text-wrap">{ contact.direct_message.last_message.message_data.slice(0, 70) + (contact.direct_message.last_message.message_data.length > 70 ? '...' : '') }</Text>
                      <Text className="ml-auto">{ formatTimestamp(contact.direct_message.last_message.timestamp) }</Text>
                    </>
                  }
              </VStack>
            </HStack>
          </VStack>
        </Card>
      </Pressable>
    )
  }