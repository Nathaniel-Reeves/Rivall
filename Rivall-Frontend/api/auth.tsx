const API_URL = process.env.EXPO_PUBLIC_API_URL;

export async function getUser() {
    const endpoint = `${API_URL}/api/v1/users?username=test_username`;
    
    console.log('fetching:', endpoint);
    return fetch(endpoint)
    .then(response => {
      if (!response.ok) {
        throw new Error('Network response was not ok');
      }
      return response.json().then(data => {
          console.log(data);
          return data;
      });
    })
    .catch((error) => {
        console.error('Error:', error);
    });
}

export async function login(email: string, password: string) {
  const endpoint = `${API_URL}/api/v1/auth/login`;
  console.log('fetching:', endpoint);
  const res = await fetch(endpoint, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ email, password }),
  })

  const data = await res.json();
  if (!res.ok) {
    throw new Error(data.message);
  }
  return data;
}

export async function register(email: string, password: string) {
  const endpoint = `${API_URL}/api/v1/auth/register`;
  console.log('fetching:', endpoint);
  const res = await fetch(endpoint, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ email, password }),
  })

  const data = await res.json();
  if (!res.ok) {
    throw new Error(data.message);
  }
  return data;
}
