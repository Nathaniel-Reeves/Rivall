import { Stack } from 'expo-router';
import "@/global.css";
import { GluestackUIProvider } from '@/components/ui/gluestack-ui-provider';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';

import { BackgroundGradientWrapper } from '@/components/BackgroundGradientWrapper';

const queryClient = new QueryClient();

export default function RootLayout() {
  return (
    <QueryClientProvider client={queryClient}>
      <GluestackUIProvider>
        <BackgroundGradientWrapper>
          <Stack screenOptions={{
            headerShown: false,
            contentStyle: { backgroundColor: 'transparent' }, // Ensure the background is transparent
          }}>
            <Stack.Screen name="index" options={{ title: 'Welcome' }} />
          </Stack>
        </BackgroundGradientWrapper>
      </GluestackUIProvider>
    </QueryClientProvider>
  )
}
{/* <GluestackUIProvider>
        <Stack screenOptions={{
            headerLeft: () => (
              <Link href={"(auth)/login"} asChild>
                <Pressable className="flex-row gap-2">
                  <Icon as={User} />
                </Pressable>
              </Link>
          ),
          headerTitleAlign: 'center',
        }}></Stack> */}