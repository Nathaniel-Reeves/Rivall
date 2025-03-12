export async function login(email: string, password: string) : Promise<[any, boolean]> {
  const endpoint = `${process.env.EXPO_PUBLIC_API_URL}/api/v1/auth/login`;
  const method = 'POST';
  const headers = {};
  const body = {
    email: email,
    password: password,
  };
  try {
    const res = await fetch(endpoint, {
      method: method,
      headers: headers,
      body: JSON.stringify(body),
    });
    const data = await res.json();
    return [data, res.ok]; 
  } catch (error: any) {
    console.error(`Error: ${error}`);
    return [null, false];
  }
}

export async function register(firstName: string, lastName: string, email: string, password: string) : Promise<[any, boolean]> {
  const endpoint = `${process.env.EXPO_PUBLIC_API_URL}/api/v1/auth/register`;
  const method = 'POST';
  const headers = {};
  const body = {
    email: email,
    password: password,
    first_name: firstName,
    last_name: lastName,
  };
  try {
    const res = await fetch(endpoint, {
      method: method,
      headers: headers,
      body: JSON.stringify(body),
    });
    const data = await res.json();
    return [data, res.ok];
  } catch (error: any) {
    console.error(`Error: ${error}`);
    return [null, false];
  }
}

export async function sendCodeToEmail(email: string) : Promise<[any, boolean]> {
  const endpoint = `${process.env.EXPO_PUBLIC_API_URL}/api/v1/auth/recovery/send-code`;
  const method = 'POST';
  const headers = {};
  const body = {
    email: email,
  };
  try {
    const res = await fetch(endpoint, {
      method: method,
      headers: headers,
      body: JSON.stringify(body),
    });
    const data = await res.json();
    return [data, res.ok];
  } catch (error: any) {
    console.error(`Error: ${error}`);
    return [null, false];
  }
}
