export async function getUser(id: string, token: string) : Promise<[any, boolean]> {
  const endpoint = `${process.env.EXPO_PUBLIC_API_URL}/api/v1/users/${id}`;
  const method = 'GET';
  const headers = { 'Authorization': token };
  // const body = {};
  try {
    const res = await fetch(endpoint, {
      method: method,
      headers: headers,
      // body: JSON.stringify(body),
    });
    if (!res.ok) {
      return [null, false];
    }
    const data = await res.json();
    return [data, true];
  } catch (error) {
    console.error('Error:', error);
    return [null, false];
  }
}
