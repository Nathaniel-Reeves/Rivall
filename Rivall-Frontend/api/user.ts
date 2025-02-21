import { wrappedFetch } from './wrappedFetch';

export async function getUser(id: string, token: string) : Promise<Object> {
    const endpoint = `${process.env.EXPO_PUBLIC_API_URL}/api/v1/users/${id}`;
    const method = 'GET';
    const headers = { 'Authorization': token };
    const body = {};
    return await wrappedFetch(endpoint, method, headers, body);
}
