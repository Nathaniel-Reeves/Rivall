import { Image } from '@/components/ui/image';
import { Box } from '@/components/ui/box';
import { VStack } from '@/components/ui/vstack';
import { Button } from '@/components/ui/button';
import { Text } from '@/components/ui/text';
import { Link } from 'expo-router';

export default function WelcomeScreen() {
  return (
    <Box className="flex-1 justify-center w-80 mx-auto">
      <Image
        source={require('@/assets/icon.png')}
        className="shadow-lg shadow-black w-[236px] h-[236px] justify-center mx-auto mb-10 rounded-[52px]"
      />
      <Text className="text-typography-800 text-2xl font-medium text-pretty text-center mb-20">Where Rivalls' Become Campions</Text>
      <VStack className="gap-4 mt-10">
        <Link href="/login" asChild replace>
          <Button
            className="shadow-md shadow-black"
            size="lg"
          >
            <Text className="text-typography-0 text-lg">Login</Text>
          </Button>
        </Link>
        <Link href="/registration" asChild replace>
          <Button
            className="shadow-md shadow-black bg-secondary-500"
            size="lg"
          >
            <Text className="text-typography-800 text-lg">Register</Text>
          </Button>
        </Link>
        <Text className="text-center text-lg font-medium">Compete! Win! Brag!</Text>
      </VStack>
    </Box>
  );
}
