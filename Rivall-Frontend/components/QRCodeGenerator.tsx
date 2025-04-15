import QRCode from 'react-native-qrcode-svg';
import { Box } from '@/components/ui/box';
import { View, Text } from 'react-native';

export default function QRGenerator(props: {
  className?: string;
  value: string | undefined | null;
  size?: number;
  logo: any;
  logoSize?: number;
  logoBackgroundColor?: string;
  logoMargin?: number;
}) {
  if (!props.value) {
    return (
      <Box className={props.className ? props.className : ''}>
        <View style={{
          minWidth: props.size ? props.size : 250,
          minHeight: props.size ? props.size : 250,
          alignItems: 'center',
        }}>
          <Text>No Data</Text>
        </View>
      </Box>
    )
  }
  return (
    <Box className={props.className ? props.className : ''}>
      <View style={{
        minWidth: props.size ? props.size : 250,
        minHeight: props.size ? props.size : 250,
        alignItems: 'center',
      }}>
        <QRCode
          size={props.size ? props.size : 250}
          value={props.value ? props.value : ''}
          backgroundColor={props.logoBackgroundColor ? props.logoBackgroundColor : 'transparent'}
          logo={props.logo}
          logoSize={props.logoSize ? props.logoSize : 30}
          logoMargin={props.logoMargin ? props.logoMargin : 0}
          logoBackgroundColor={props.logoBackgroundColor ? props.logoBackgroundColor : 'transparent'}
        />
      </View>
    </Box>
  )
};