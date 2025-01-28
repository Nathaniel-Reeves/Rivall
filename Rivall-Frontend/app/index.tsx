import { StatusBar } from 'expo-status-bar';
import { StyleSheet, Text, View, FlatList } from 'react-native';
import products from '../assets/products.json';

import ProductListItem from '../components/ProductListItem';

export default function HomeScreen() {
  return (
    <FlatList data={products} renderItem={({ item }) => (<ProductListItem product={item} />)} />
  );
}