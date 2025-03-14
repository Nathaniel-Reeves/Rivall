import { create } from 'zustand';

interface User {
  _id: string;
  token: string;
  first_name: string;
  last_name: string;
  email: string;
  avatar_image: string;
}

type State = {
  otp: string;
  access_token: string;
  access_token_expires_at: string;
  refresh_token: string;
  refresh_token_expires_at: string;
  user: User;
}

type Actions = {
  setUserData: (userData: Partial<User>) => void;
  setAuthData: (authData: Partial<State>) => void;
  setStoreData: (data: Partial<State>) => void;
  clearStore: () => void;
}

export const useUserStore = create<State & Actions>()(
  (set) => ({
    otp: '',
    access_token: '',
    access_token_expires_at: '',
    refresh_token: '',
    refresh_token_expires_at: '',
    user: {
      _id: '',
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
    setAuthData: (authData: Partial<State>) => set({
      otp: authData.otp,
      access_token: authData.access_token,
      access_token_expires_at: authData.access_token_expires_at,
      refresh_token: authData.refresh_token,
      refresh_token_expires_at: authData.refresh_token_expires_at
    }),
    setStoreData: (data: Partial<State>) => set({
      otp: data.otp,
      access_token: data.access_token,
      access_token_expires_at: data.access_token_expires_at,
      refresh_token: data.refresh_token,
      refresh_token_expires_at: data.refresh_token_expires_at,
      user: { 
        _id: data.user?._id || '',
        token: data.user?.token || '',
        first_name: data.user?.first_name || '',
        last_name: data.user?.last_name || '',
        email: data.user?.email || '',
        avatar_image: data.user?.avatar_image || '',
       }
    }),
    clearStore: () => {
      set({ 
        otp: '',
        access_token: '',
        access_token_expires_at: '',
        refresh_token: '',
        refresh_token_expires_at: '',
        user: {
          _id: '',
          token: '',
          first_name: '',
          last_name: '',
          email: '',
          avatar_image: '',
        },
      })
    },
  })
);