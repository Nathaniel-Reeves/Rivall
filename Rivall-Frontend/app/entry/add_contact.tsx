import { BackgroundGradientWrapper } from '@/components/BackgroundGradientWrapper';
import QRScanner from '@/components/QRCodeScanner';
import { Card } from '@/components/ui/card';
import { Heading } from '@/components/ui/heading';
import { Box } from '@/components/ui/box';
import { Avatar, AvatarFallbackText, AvatarImage } from '@/components/ui/avatar';
import { Icon } from '@/components/ui/icon';
import { Text } from '@/components/ui/text';
import { Spinner } from '@/components/ui/spinner';
import { HStack } from '@/components/ui/hstack';
import { VStack } from '@/components/ui/vstack';
import { Button, ButtonIcon, ButtonText } from '@/components/ui/button';
import { ChevronLeftIcon } from '@/components/ui/icon';
import { Link, useRouter } from 'expo-router';
import { useState } from 'react';
import { getContact } from '@/api/contact';
import { useUserStore } from '@/global-store/user_store';
import { Plus, Mail, Crown } from 'lucide-react-native';
import { Contact } from '@/types';
import { addContact } from '@/api/contact';

export default function AddContactScreen() {

  const router = useRouter();
  const [lockScanner, setLockScanner] = useState(false);
  const [loading, setLoading] = useState(false);
  const [hasContact, setHasContact] = useState(false);
  const [contact, setContact] = useState<Contact>();
  const access_token = useUserStore((state: any) => state.access_token);
  const user = useUserStore((state: any) => state.user);

  const handleQRCode = async (scanned_id: string) => {
    setLockScanner(true);
    setLoading(true);
    const res = await getContact(scanned_id, access_token);
    if (res.status == 200) {
      setContact(res.data);
      setHasContact(true);
      console.log(res.data)
    } else {
      console.error('Error getting contact');
      console.error(res);
    }
    setLoading(false);
  }

  const [ error , setError ] = useState<string | null>(null);

  const handleAddContact = async () => {
    setLoading(true);
    const res = await addContact(user._id, contact._id, access_token);
    if (res.status == 201) {
      router.push('/entry/home/contacts');
    } else {
      console.log(`Error adding contact: ${res.status}`);
      setError('Error adding contact');
    }
    setLoading(false);
  }

  return (
    <BackgroundGradientWrapper>
      <Card className="align-top m-4 mt-16 p-10 shadow-md shadow-black">
        <HStack className="gap-4 mb-10">
          <Link href="/entry/home/contacts" asChild replace>
            <Button
              className="mb-6 w-20"
              variant="outline"
            >
              <ButtonIcon as={ChevronLeftIcon}/>
            </Button>
          </Link>
          <Heading className="text-typography-800 text-center">
            Scan Rivall QR Code
          </Heading>
        </HStack>
        {lockScanner || loading ? null : 
        <QRScanner
          onQRCode={(data) => handleQRCode(data)}
        />}
        {loading ? 
          <HStack className="gap-2 justify-center">
            <Spinner size="large" />
          </HStack>
          : null
        }
        {contact && !loading ? 
          <VStack className="mx-auto gap-4 w-full">
            <Box className="w-full">
              <Avatar className={`w-28 h-28 mr-4 my-auto mx-auto`} style={{ backgroundColor: "blue" }}>
                <AvatarFallbackText className="text-white text-4xl">
                  {contact.first_name[0] + ' ' + contact.last_name[0]}
                </AvatarFallbackText>
                <AvatarImage
                  source={{ uri: contact.avatar_image }}
                  className="w-28 h-28 rounded-full"
                />
              </Avatar>
            </Box>
            <HStack className="gap-2">
              <Icon as={Crown} size="lg" className="my-auto text-black"/>
              <Text className="text-typography-800 text-lg text-pretty text-left my-auto">
                {contact.first_name} {contact.last_name}
              </Text>
            </HStack>
            <HStack className="gap-2">
              <Icon as={Mail} size="lg" className="my-auto"/>
              <Text className="text-typography-800 text-lg text-pretty text-left my-auto">
                {contact.email}
              </Text>
            </HStack>
            <Button onPress={handleAddContact}>
              <ButtonIcon as={Plus}/>
              <ButtonText>
                Add { contact.first_name } to Contacts
              </ButtonText>
            </Button>
          </VStack>
          : null
        }
        {}
      </Card>
    </BackgroundGradientWrapper>
  )
}