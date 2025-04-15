import { client } from './axios.config';

export async function getUser(id: string, token: string) : Promise<any> {
  const endpoint = `/users/${id}`;
  const headers = { 'Authorization': token };
  // const body = {};
  try {
    const res = await client.get(endpoint, { headers: headers });
    return res;
  } catch (error) {
    console.error('Error:', error);
    return {status: error.status, data: null};
  }
}
