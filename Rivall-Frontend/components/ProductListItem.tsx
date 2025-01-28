import { Text } from 'react-native'

export default function ProductListItem ({ product }: { product: { name: string } }) {
    return (
      <Text>{product.name}</Text>
    )
  }