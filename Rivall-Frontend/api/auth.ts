import { wrappedFetch } from './wrappedFetch';

export async function login(email: string, password: string) : Promise<Object> {
  const endpoint = `${process.env.EXPO_PUBLIC_API_URL}/api/v1/auth/login`;
  const method = 'POST';
  const headers = {};
  const body = {
    email: email,
    password: password,
  };
  return await wrappedFetch(endpoint, method, headers, body);
}

export async function register(email: string, password: string) : Promise<Object> {
  const endpoint = `${process.env.EXPO_PUBLIC_API_URL}/api/v1/auth/register`;
  const method = 'POST';
  const headers = {};
  const body = {
    email: email,
    password: password,
  };
  return await wrappedFetch(endpoint, method, headers, body);
}
