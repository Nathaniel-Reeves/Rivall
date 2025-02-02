import { Stack, Link } from 'expo-router';
import "@/global.css";
import { GluestackUIProvider } from '@/components/ui/gluestack-ui-provider';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { Icon } from '@/components/ui/icon';
import { Pressable } from 'react-native';
import { User } from 'lucide-react-native';

export default function RootLayout() {
  return (
    <QueryClientProvider client={new QueryClient()}>
      <GluestackUIProvider>
        <Stack screenOptions={{
            headerLeft: () => (
              <Link href={"(auth)/login"} asChild>
                <Pressable className="flex-row gap-2">
                  <Icon as={User} />
                </Pressable>
              </Link>
          ),
          headerTitleAlign: 'center',
        }}>
          <Stack.Screen name="index" options={{ title: 'Shop' }} />
          <Stack.Screen name="product/[id]" options={{ title: 'Product Details' }} />
        </Stack>
      </GluestackUIProvider>
    </QueryClientProvider>
  )
}