import { Image } from '@/components/ui/image';
import { Box } from '@/components/ui/box';
import { Text } from '@/components/ui/text';

export default function StartupErrorScreen() {
  return (
    <Box className="flex-1 justify-center w-80 mx-auto">
      <Image
        source={require('@/assets/icon.png')}
        className="shadow-md shadow-black w-[236px] h-[236px] justify-center mx-auto mb-10 rounded-[42px]"
        alt="Rivall Logo"
      />
      <Text className="text-typography-800 text-2xl font-medium text-pretty text-center mb-20">Where Rivalls' Become Campions</Text>
      <Text className="text-typography-800 text-2xl font-medium text-pretty text-center mb-20">Startup Error</Text>
      <Text className="text-typography-800 text-2xl font-medium text-pretty text-center mb-20">Please check your internet connection and try again</Text>
    </Box>
  )
}
