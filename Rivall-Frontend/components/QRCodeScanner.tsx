import { Card } from "@/components/ui/card"
import { Button, ButtonText } from '@/components/ui/button';
import { Text } from '@/components/ui/text';

import { CameraView, useCameraPermissions } from 'expo-camera';
import { View } from 'react-native';
import { useState } from 'react';

export default function QRScanner(props: {
  onQRCode?: (data: string) => void;
}) {
  const [permission, requestPermission] = useCameraPermissions();
  const [lockScanner, setLockScanner] = useState(false);

  if (!permission) {
    // Camera permissions are still loading.
    return <View />;
  }

  if (!permission.granted) {
    // Camera permissions are not granted yet.
    return (
      <Card
        className="block min-h-[400px] my-2 py-2"
      >
        <Text>We need your permission use the scanner.</Text>
        <Button onPress={requestPermission}>
          <ButtonText>Grant Permission</ButtonText>
        </Button>
      </Card>
    );
  }

  const handleQRCode = (data: string) => {
    if (props.onQRCode) {
      props.onQRCode(data);
    }
  }

  return (
    <Card
      className="block min-h-[400px] my-2 py-2"
    >
      <CameraView
        style={{
          flex: 1,
          zIndex: 100000
        }}
        facing='back'
        onBarcodeScanned={({ data }) => {
          if (data && !lockScanner) {
            setLockScanner(true);
            handleQRCode(data);
            setTimeout(async () => {
              setLockScanner(false);
            }, 2000);
          }
        }}
      >
      </CameraView>
    </Card>
  );
}