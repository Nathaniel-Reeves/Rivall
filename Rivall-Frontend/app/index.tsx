import { useEffect, useState } from 'react';
import { useQuery } from '@tanstack/react-query';

import { LinearGradient } from 'expo-linear-gradient';
import { Image } from '@/components/ui/image';
import { Box } from '@/components/ui/box';
import { Text } from '@/components/ui/text';
import { Spinner } from '@/components/ui/spinner';
import colors from "tailwindcss/colors"
import WelcomeScreen from './welcome';

function HomeScreen() {
  return (
    <LinearGradient
      // Vertical Background Linear Gradient
      colors={['#77FBFF', '#26C1FE']}
      start={[0, 0]}
      end={[0, 1]}
      style={{
        position: 'absolute',
        left: 0,
        right: 0,
        top: 0,
        height: '100%',
      }}
    >
      <Box className="flex-1 justify-center w-80 mx-auto">
        <Image
          source={require('@/assets/icon.png')}
          className="shadow-md shadow-black w-[236px] h-[236px] justify-center mx-auto mb-10 rounded-[42px]"
          alt="Rivall Logo"
        />
        <Text className="text-typography-800 text-2xl font-medium text-pretty text-center mb-20">Where Rivalls' Become Campions</Text>
        <Spinner size="large" color={colors.gray[700]} />
      </Box>
    </LinearGradient>
  )
}

export default function App() {

  const [ isLoggedIn, setIsLoggedIn ] = useState<boolean | null>(null);

  useEffect(() => {
    const timeoutId = setTimeout(() => {
      setIsLoggedIn(false);
    }, 5000);
    return () => clearTimeout(timeoutId);
  }, []);

  // const { data, isLoading, error } = useQuery({
  //   queryKey: ['user'],
  //   queryFn: getUser
  // });

  return (
    <>
      {isLoggedIn == false ? <WelcomeScreen/> : <HomeScreen/>}
    </>
  )
}