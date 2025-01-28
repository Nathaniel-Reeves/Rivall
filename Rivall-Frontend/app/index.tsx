import { StatusBar } from 'expo-status-bar';
import { StyleSheet, Text, View, FlatList, useWindowDimensions } from 'react-native';
import products from '@/assets/products.json';
import { useBreakpointValue } from '@/components/ui/utils/use-break-point-value';

import ProductListItem from '@/components/ProductListItem';

export default function HomeScreen() {
  // const { width } = useWindowDimensions();
  // const numColumns = width > 768 ? 3 : 2; // re-renders on screen width change

  const numColumns = useBreakpointValue({ // uses breakpoints from tailwind.config.js
    default: 2,
    sm: 3,
    md: 4,
  });

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