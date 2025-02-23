import { Image } from '@/components/ui/image';
import { Box } from '@/components/ui/box';
import { Text } from '@/components/ui/text';
import { BackgroundGradientWrapper } from '@/components/BackgroundGradientWrapper';
import { useUserStore } from '@/global-store/user_store';

export default function AccountScreen() {

  const user = useUserStore((state: any) => state.user);

  return (
    <BackgroundGradientWrapper>
      <Box className="flex-1 justify-center w-80 mx-auto">
        <Image
          source={require('@/assets/icon.png')}
          className="shadow-md shadow-black w-[236px] h-[236px] justify-center mx-auto mb-10 rounded-[42px]"
          alt="Rivall Logo"
        />
        <Text className="text-typography-800 text-2xl font-medium text-pretty text-center mb-20">
          first name:{user.first_name}
        </Text>
        <Text className="text-typography-800 text-2xl font-medium text-pretty text-center mb-20">
          last name:{user.last_name}
        </Text>
        <Text className="text-typography-800 text-2xl font-medium text-pretty text-center mb-20">
          user:{JSON.stringify(user)}
        </Text>
      </Box>
    </BackgroundGradientWrapper>
  )
}