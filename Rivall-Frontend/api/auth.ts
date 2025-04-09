import { client } from './axios.config';

export async function login(email: string, password: string) : Promise<any> {
  const endpoint = '/auth/login';
  const body = {
    email: email,
    password: password,
  };
  try {
    const res = await client.post(endpoint, body);
    return res; 
  } catch (error: any) {
    console.error(`Error: ${error}`);
    return null;
  }
}

export async function register(firstName: string, lastName: string, email: string, password: string) : Promise<any> {
  const endpoint = '/auth/register';
  const body = {
    email: email,
    password: password,
    first_name: firstName,
    last_name: lastName,
  };
  try {
    const res = await client.post(endpoint, body);
    return res;
  } catch (error: any) {
    console.error(`Error: ${error}`);
    return null;
  }
}

export async function sendCodeToEmail(email: string) : Promise<any> {
  const endpoint = '/auth/recovery/send-code';
  const body = {
    email: email,
  };
  try {
    const res = await client.post(endpoint, body);
    return res;
  } catch (error: any) {
    console.error(`Error: ${error}`);
    return null;
  }
}

export async function validateAccountRecoveryCode(email: string, code: string) : Promise<any> {
  const endpoint = '/auth/recovery/validate-code';
  const body = {
    email: email,
    code: code,
  };
  try {
    const res = await client.post(endpoint, body);
    return res;
  } catch (error: any) {
    console.error(`Error: ${error}`);
    return null;
  }
}
export async function resetPassword(state: any, password: string) : Promise<any> {
  const endpoint = `/auth/recovery/${state.user._id}/reset-password`;
  const headers = {
    'Authorization': state.access_token,
  };
  const body = {
    _id: state.user._id,
    password: password,
  };
  try {
    const res = await client.put(endpoint, body, { headers: headers });
    return res;
  } catch (error: any) {
    console.error(`Error: ${error}`);
    return null;
  }
}
