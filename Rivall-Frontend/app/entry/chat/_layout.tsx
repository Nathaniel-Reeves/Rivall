import { Stack } from 'expo-router';

export default function ChatScreenLayout() {
  return (
    <Stack screenOptions={{
      headerShown: true,
      statusBarTranslucent: true,
      statusBarBackgroundColor: 'transparent',
      statusBarStyle: 'dark',
      contentStyle: { backgroundColor: 'transparent' }, // Ensure the background is transparent
    }}>
    </Stack>
  )
}