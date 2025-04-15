import { useQuery } from '@tanstack/react-query';
import { Image } from '@/components/ui/image';
import { Box } from '@/components/ui/box';
import { Text } from '@/components/ui/text';
import { Spinner } from '@/components/ui/spinner';
import colors from "tailwindcss/colors"
import { Redirect } from 'expo-router';
import WelcomeScreen from './(auth)/welcome';
import StartupErrorScreen from './(auth)/startup_error';
import { getUser } from '@/api/user';
import { useUserStore } from '@/global-store/user_store';

function HomeScreen() {
  return (
    <Box className="flex-1 justify-center w-80 mx-auto">
      <Image
        source={require('@/assets/icon.png')}
        className="shadow-md shadow-black w-[236px] h-[236px] justify-center mx-auto mb-10 rounded-[42px]"
        alt="Rivall Logo"
      />
      <Text className="text-typography-800 text-2xl font-medium text-pretty text-center mb-20">Where Rivalls' Become Campions</Text>
      <Spinner size="large" color={colors.gray[700]} />
    </Box>
  )
}

export default function App() {

  // Check if storage has user id and token
  const user = useUserStore((state: any) => state.user);
  const access_token = useUserStore((state: any) => state.access_token);
  if (user._id == '') {
    console.debug('No user id')
    return (
      <WelcomeScreen/> // Redirect user to login or register
    )
  }

  // Get User Data using auth token
  const { data, isLoading, error } = useQuery({
    queryKey: ['getUser', 'Startup'],
    queryFn: () => getUser(user._id, access_token),
    retryDelay: attempt => Math.min(attempt > 1 ? 2 ** attempt * 1000 : 1000, 30 * 1000),
  });

  if (isLoading) {
    return (
      <HomeScreen/>
    )
  }

  if (error) {
    // TODO: make login error screen
    console.error(error)
    return (
      <StartupErrorScreen/>
    )
  }

  if (data?._id == undefined) {
    return (
      <WelcomeScreen/> // Redirect user to login or register
    )
  }

  // Redirect user to home screen
  return (
    <Redirect href="/entry/home"/>
  )
}