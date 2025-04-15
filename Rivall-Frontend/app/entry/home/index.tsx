import { Box } from '@/components/ui/box';
import { Card } from '@/components/ui/card';
import { Text } from '@/components/ui/text';
import { Icon } from '@/components/ui/icon';
import { Button, ButtonText } from '@/components/ui/button';
import { Heading } from '@/components/ui/heading';
import { Avatar, AvatarFallbackText, AvatarImage } from '@/components/ui/avatar';
import { BackgroundGradientWrapper } from '@/components/BackgroundGradientWrapper';
import { useUserStore } from '@/global-store/user_store';
import { ScrollView } from 'react-native';
import { useRouter } from 'expo-router';
import { HStack } from '@/components/ui/hstack';
import { Camera, Mail, Crown } from 'lucide-react-native';
import QRGenerator from '@/components/QRCodeGenerator';

export default function AccountScreen() {

  // Check if storage has user id and token
  const user = useUserStore((state: any) => state.user);
  const logout = useUserStore((state: any) => state.clearStore);
  const router = useRouter();

  let logo = require('@/assets/rivall-logos/black_shape_white_r.png');

  const handleLogout = async () => {
    // Clear user data from store
    console.debug('User logged out')
    logout();
    router.replace('/(auth)/welcome');
  }

  return (
    <BackgroundGradientWrapper>
      <ScrollView className="h-full">
        <Box className="flex-1 justify-center m-4">
          <HStack className="gap-3">
            <Box className="w-24 h-24 shadow-lg shadow-black mx-auto my-auto rounded-full">
              <Avatar className="w-24 h-24 bg-white">
                {/* <Icon as={Camera} size="lg" className="my-auto w-10 h-10 text-black"/> */}
                <AvatarFallbackText className="text-black text-2xl">
                  {user.first_name[0] + ' ' + user.last_name[0]}
                </AvatarFallbackText>
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
        <Card className="my-auto mx-4 mb-4 gap-4">
          <Heading className="text-typography-800 text-center text-2xg">
            Scan With Your Rivall App
          </Heading>
          <QRGenerator
            value={user._id}
            size={300}
            logo={ logo }
            logoSize={80}
            logoBackgroundColor='white'
            logoMargin={2}
          />
          <Text className="text-typography-800 text-center text-md">
            User ID: {user._id}
          </Text>
          <Button onPress={handleLogout}>
            <ButtonText>Logout</ButtonText>
          </Button>
        </Card>
      </ScrollView>
    </BackgroundGradientWrapper>
  )
}