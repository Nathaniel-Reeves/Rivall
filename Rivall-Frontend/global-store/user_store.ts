import { create } from 'zustand';

interface User {
  id: string;
  token: string;
  first_name: string;
  last_name: string;
  email: string;
  avatar_image: string;
}

type State = {
  user: User;
}

type Actions = {
  setUserData: (userData: Partial<User>) => void;
  clearUserData: () => void;
}

export const useUserStore = create<State & Actions>()(
  (set) => ({
    user: {
      id: '',
      token: '',
      first_name: '',
      last_name: '',
      email: '',
      avatar_image: '',
    },
    setUserData: (refreshUser: Partial<User>) => set({
      user: { 
        id: refreshUser.id || '',
        token: refreshUser.token || '',
        first_name: refreshUser.first_name || '',
        last_name: refreshUser.last_name || '',
        email: refreshUser.email || '',
        avatar_image: refreshUser.avatar_image || '',
       }
    }),
    clearUserData: () => set({ 
      user: {
        id: '',
        token: '',
        first_name: '',
        last_name: '',
        email: '',
        avatar_image: '',
      }
    }),
  })
);