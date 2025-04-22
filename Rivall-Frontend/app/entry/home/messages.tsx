import { Fab, FabIcon } from '@/components/ui/fab';
import { BackgroundGradientWrapper } from '@/components/BackgroundGradientWrapper';
import { FlatList } from 'react-native';
import { useRouter } from 'expo-router';
import { useQuery } from '@tanstack/react-query';
import { Plus } from 'lucide-react-native';
import { getUser } from '@/api/user';
import { useUserStore } from '@/global-store/user_store';
import ContactCard from '@/components/ContactCard';

export default function ContactsScreen() {
  const user = useUserStore((state: any) => state.user);
  const access_token = useUserStore((state: any) => state.access_token);

  // Get User Data using auth token
  const { data, isLoading, error } = useQuery({
    queryKey: ['getUser', 'Startup'],
    queryFn: () => getUser(user._id, access_token),
    retryDelay: attempt => Math.min(attempt > 1 ? 2 ** attempt * 1000 : 1000, 30 * 1000),
  });

  console.log(JSON.stringify(data, null, 2));
  const contacts = data?.data?.populated_contacts || [];
  const router = useRouter();

  return (
    <BackgroundGradientWrapper>
      <Fab className="bg-sky-800" onPress={() => router.push('/entry/add_contact')}>
        <FabIcon as={Plus} size="2xl" className="text-white"/>
      </Fab>
      <FlatList
        data={contacts}
        renderItem={({ item }) => (
          <ContactCard key={item._id} contact={item} />
        )}
      >
      </FlatList>
    </BackgroundGradientWrapper>
  )
}
