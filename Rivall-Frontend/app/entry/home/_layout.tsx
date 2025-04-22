import React from 'react';
import { Tabs } from 'expo-router';
import { Icon } from "@/components/ui/icon"
import { CircleUserRound, MessageSquare, Share2 } from 'lucide-react-native';
import { Button, ButtonIcon } from '@/components/ui/button';
import { HStack } from '@/components/ui/hstack';
import * as Linking from 'expo-linking';
import { Share } from 'react-native';
import { useUserStore } from '@/global-store/user_store';

const Layout = () => {

  const user = useUserStore((state: any) => state.user);

  const handleShare = async () => {
    try {
      const url = Linking.createURL('entry/home?user=' + user._id)
      const result = await Share.share({
        message: `Add me as a contact on Rivall!\n${url}`,
        url: url,
      });
      if (result.action === Share.sharedAction) {
        console.log('Shared successfully');
      } else if (result.action === Share.dismissedAction) {
        console.log('Share dismissed');
      }
    } catch (error) {
      console.error('Error sharing:', error);
    }
  }

  return (
    <Tabs screenOptions={{ tabBarActiveTintColor: '#1A8BC1' }}>
      <Tabs.Screen
      name="index"
      options={{
        title: 'Account',
        headerShown: true,
        tabBarIcon: ({ color }) => <Icon size='xl' color={color} as={CircleUserRound} />,
        headerRight: () => (<HStack className="gap-6 m-4">
          <Button onPress={handleShare} variant="link" size="lg" className="rounded-full"><ButtonIcon as={Share2} size="lg" className="w-[28px] h-[28px]"/></Button>
        </HStack>),
      }}
      />
      <Tabs.Screen
      name="messages"
      options={{
        title: 'Messages',
        headerShown: true,
        tabBarIcon: ({ color }) => <Icon size='xl' color={color} as={MessageSquare} />,
      }}
      />
    </Tabs>
  );
};

export default Layout;