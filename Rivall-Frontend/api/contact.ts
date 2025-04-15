import { client } from './axios.config';

export async function getContact(id: string, token: string) : Promise<any> {
  const endpoint = `/contacts/${id}`;
  const headers = { 'Authorization': token };
  try {
    const res = await client.get(endpoint, { headers: headers });
    return res;
  } catch (error) {
    console.error('Error:', error);
    return {status: error.status, data: null};
  }
}

export async function getChat(id: string, token: string, chatID: string) : Promise<any> {
  const endpoint = `/users/${id}/contacts/${chatID}/chat`;
  const headers = { 'Authorization': token };
  try {
    const res = await client.get(endpoint, { headers: headers });
    return res;
  } catch (error) {
    console.error('Error:', error);
    return {status: error.status, data: null};
  }
}

export async function addContact(user_id: string, contact_id:string, token: string) : Promise<any> {
  const endpoint = `/users/${user_id}/contacts`;
  const headers = { 'Authorization': token };
  const body = {
    contact_id: contact_id,
  };
  try {
    const res = await client.post(endpoint, body, { headers: headers });
    return res;
  } catch (error: any) {
    console.error('Error:', JSON.stringify(error, null, 2));
    return {status: error.status, data: null};
  }
}