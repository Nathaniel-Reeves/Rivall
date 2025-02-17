import { Stack } from 'expo-router';
import "@/global.css";
import { GluestackUIProvider } from '@/components/ui/gluestack-ui-provider';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';

export default function RootLayout() {
  return (
    <QueryClientProvider client={new QueryClient()}>
      <GluestackUIProvider>
        <Stack screenOptions={{
          headerShown: false,
        }}>
          <Stack.Screen name="index" options={{ title: 'Welcome' }} />
          <Stack.Screen name="login" options={{ title: 'Login' }} />
          <Stack.Screen name="registration" options={{ title: 'Register' }} />
          <Stack.Screen name="product/[id]" options={{ title: 'Product Details' }} />
        </Stack>
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