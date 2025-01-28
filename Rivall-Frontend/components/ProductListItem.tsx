import { View, Text } from 'react-native'
import { Button, ButtonText } from '@/components/ui/button'

export default function ProductListItem ({ product }: { product: { name: string } }) {
    return (
      <View>
        <Text>{product.name}</Text>
        <Button>
          <ButtonText>Add to cart</ButtonText>
        </Button>
      </View>
    )
  }