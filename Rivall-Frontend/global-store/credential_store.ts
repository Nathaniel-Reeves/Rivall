import { create } from 'zustand';
import { persist, createJSONStorage } from 'zustand/middleware';

export const useCredentials = create(
  persist(
    (set) => ({
      credentials: {
        userID: '',
        token: '',
      },
      setCredentials: (userID: string, token: string) => set({ credentials: { userID, token } }),
      clearCredentials: () => set({ credentials: { userID: '', token: '' } }),
    }),
    {
      name: 'credentials-storage', // unique name
      storage: createJSONStorage(() => localStorage), // use localStorage
    }
  )
);