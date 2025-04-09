import { Image } from '@/components/ui/image';
import { Box } from '@/components/ui/box';
import { Card } from '@/components/ui/card';
import { Text } from '@/components/ui/text';
import { Icon } from '@/components/ui/icon';
import { Avatar, AvatarFallbackText, AvatarImage } from '@/components/ui/avatar';
import { BackgroundGradientWrapper } from '@/components/BackgroundGradientWrapper';
import { useUserStore } from '@/global-store/user_store';
import { ScrollView } from 'react-native';
import { HStack } from '@/components/ui/hstack';
import { Camera, Mail, Crown } from 'lucide-react-native';
import { RivallBlackCrown } from '@/components/RivallBlackCrown';

export default function AccountScreen() {

  const user = useUserStore((state: any) => state.user);

  return (
    <BackgroundGradientWrapper>
      <ScrollView className="h-full">
        <Box className="flex-1 justify-center m-4">
          <HStack className="gap-3">
            <Box className="w-24 h-24 shadow-lg shadow-black mx-auto my-auto rounded-full">
              <Avatar className="w-24 h-24 bg-white">
                <Icon as={Camera} size="lg" className="my-auto w-10 h-10 text-black"/>
                <AvatarImage
                  source={{ uri: user.avatar_image }}
                  className="w-24 h-24 rounded-full"
                />
              </Avatar>
            </Box>
            <Card className="shadow-lg shadow-black w-80 mx-auto gap-4">
              <HStack className="gap-2">
                <Icon as={Crown} size="lg" className="my-auto text-black"/>
                <Text className="text-typography-800 text-sm text-pretty text-left my-auto">
                  {user.first_name} {user.last_name}
                </Text>
              </HStack>
              <HStack className="gap-2">
                <Icon as={Mail} size="lg" className="my-auto"/>
                <Text className="text-typography-800 text-sm text-pretty text-left my-auto">
                  {user.email}
                </Text>
              </HStack>
            </Card>
          </HStack>
        </Box>
        <Card className="my-auto mx-4 mb-4">
          <Text className="text-typography-800 text-md text-pretty text-left my-auto">
            {user.first_name} {user.last_name}
          </Text>
          <Text className="text-typography-800 text-md text-pretty text-left my-auto">
            {user.email}
          </Text>
        </Card>
      </ScrollView>
    </BackgroundGradientWrapper>
  )
}