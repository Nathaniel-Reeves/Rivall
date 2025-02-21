import { LinearGradient } from 'expo-linear-gradient';
import { ReactNode } from 'react';
import { View } from 'react-native';

export const BackgroundGradientWrapper = ({ children }: { children: ReactNode }) => {
  return (
    <View style={{ flex: 1 }}>
      <LinearGradient
        // Vertical Background Linear Gradient
        colors={['#77FBFF', '#26C1FE']}
        start={[0, 0]}
        end={[0, 1]}
        style={{
          position: 'absolute',
          left: 0,
          right: 0,
          top: 0,
          height: '100%',
        }}
      >
        {children}
      </LinearGradient>
    </View>
  )
}