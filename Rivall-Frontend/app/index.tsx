import { StatusBar } from 'expo-status-bar';
import { StyleSheet, Text, View, FlatList, useWindowDimensions } from 'react-native';
import { useBreakpointValue } from '@/components/ui/utils/use-break-point-value';
import { useEffect, useState } from 'react';
import { useQuery } from '@tanstack/react-query';

import ProductListItem from '@/components/ProductListItem';

import products from '@/assets/products.json';
import { getUser } from '@/api/auth';

export default function HomeScreen() {

  // const [ user, setUser ] = useState(null);

  // useEffect(() => {
  //   console.log('HomeScreen mounted');
    
  //   const fetchUser = async () => {
  //     const user = await getUser();
  //     setUser(user);
  //   };
  //   fetchUser();
  // }, []);

  const { data, isLoading, error } = useQuery({
    queryKey: ['user'],
    queryFn: getUser
  });

  // const { width } = useWindowDimensions();
  // const numColumns = width > 768 ? 3 : 2; // re-renders on screen width change

  const numColumns = useBreakpointValue({ // uses breakpoints from tailwind.config.js
    default: 2,
    sm: 3,
    md: 4,
  });

  if (isLoading) return <Text>Loading...</Text>;
  if (error) return <Text>Error: {error.message}</Text>;

  return (
    <FlatList
      key={numColumns}
      data={products}
      renderItem={({ item }) => (<ProductListItem product={item} />)}
      numColumns={numColumns}
      contentContainerClassName="gap-2 max-w-[960px] mx-auto w-full p-3"
      columnWrapperClassName="gap-2"
    />
  );
}