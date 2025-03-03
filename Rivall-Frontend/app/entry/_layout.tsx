import React from 'react';
import { Tabs } from 'expo-router';
import { Icon } from "@/components/ui/icon"
import { CircleUserRound, List } from 'lucide-react-native';

const Layout = () => {
  return (
    <Tabs screenOptions={{ tabBarActiveTintColor: '#1A8BC1' }}>
      <Tabs.Screen
        name="index"
        options={{
          title: 'Account',
          headerShown: false,
          tabBarIcon: ({ color }) => <Icon size='xl' color={color} as={CircleUserRound} />,
        }}
      />
      <Tabs.Screen
        name="contacts"
        options={{
          title: 'Contacts',
          headerShown: false,
          tabBarIcon: ({ color }) => <Icon size='xl' color={color} as={List} />,
        }}
      />
    </Tabs>
  );
};

export default Layout;