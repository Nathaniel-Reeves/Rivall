async function fetchWithTimeout(url: string, options: any): Promise<any> {
  const { timeout = 8000 } = options; // Default to 8 seconds

  /*
  TODO:
  Will throw an 'AbortError' if the request times out.
  https://medium.com/@YassineDev/how-to-timeout-a-fetch-request-2100dfee0762#:~:text=Custom%20Timeout%20for%20fetch()%20Requests&text=Here's%20a%20breakdown%3A,to%20cancel%20fetch()%20requests.
  catch (error) { error.name === 'AbortError' }

  Another possible error is a 'TypeError: Network request failed' error.
  This error comes from the api stating something went wrong on the
  back end, usually that means the api lost connection to the database.

  In any case, these errors should report a message to the user like:
  "Our servers are down, try again later."
  */

  const controller = new AbortController();
  const timer = setTimeout(() => {
    controller.abort()
    console.error(`Request Timed Out: ${url}`)
  }, timeout);

  const response = await fetch(url, {
    ...options,
    signal: controller.signal
  });
  clearTimeout(timer);

  return response;
}

export async function login(email: string, password: string) : Promise<[any, boolean]> {
  const endpoint = `${process.env.EXPO_PUBLIC_API_URL}/api/v1/auth/login`;
  const method = 'POST';
  const headers = {
    'Content-Type': 'application/json',
  };
  const body = {
    email: email,
    password: password,
  };
  try {
    const res = await fetchWithTimeout(endpoint, {
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
  const headers = {
    'Content-Type': 'application/json',
  };
  const body = {
    email: email,
    password: password,
    first_name: firstName,
    last_name: lastName,
  };
  try {
    const res = await fetchWithTimeout(endpoint, {
      method: method,
      headers: headers,
      body: JSON.stringify(body),
    });
    // Register returns no data on success
    if (res.ok) {
      return [null, true];
    } else {
      // messages are in bite strings
      const data = await res.text();
      return [data, false];
    }
  } catch (error: any) {
    console.error(`Error: ${error}`);
    return [null, false];
  }
}

export async function sendCodeToEmail(email: string) : Promise<[any, boolean]> {
  const endpoint = `${process.env.EXPO_PUBLIC_API_URL}/api/v1/auth/recovery/send-code`;
  const method = 'POST';
  const headers = {
    'Content-Type': 'application/json',
  };
  const body = {
    email: email,
  };
  try {
    const res = await fetchWithTimeout(endpoint, {
      method: method,
      headers: headers,
      body: JSON.stringify(body),
    });
    if (res.ok) {
      return [null, true];
    }
    const data = await res.text();
    return [data, res.ok];
  } catch (error: any) {
    console.error(`Error: ${error}`);
    return [null, false];
  }
}

export async function validateAccountRecoveryCode(email: string, code: string) : Promise<[any, boolean]> {
  const endpoint = `${process.env.EXPO_PUBLIC_API_URL}/api/v1/auth/recovery/validate-code`;
  const method = 'POST';
  const headers = {
    'Content-Type': 'application/json',
  };
  const body = {
    email: email,
    code: code,
  };
  try {
    const res = await fetchWithTimeout(endpoint, {
      method: method,
      headers: headers,
      body: JSON.stringify(body),
    });
    if (res.ok) {
      const data = await res.json();
      return [data, res.ok];
    } else {
      const data = await res.text();
      return [data, res.ok];
    }
  } catch (error: any) {
    console.error(`Error: ${error}`);
    return [null, false];
  }
}
export async function resetPassword(state: any, password: string) : Promise<[any, boolean]> {
  const endpoint = `${process.env.EXPO_PUBLIC_API_URL}`+encodeURI(`/api/v1/auth/recovery/${state.user._id}/reset-password`);
  console.debug(endpoint);
  const method = 'PUT';
  const headers = {
    'Content-Type': 'application/json',
    'Authorization': state.access_token,
  };
  const body = {
    _id: state.user._id,
    password: password,
  };
  try {
    const res = await fetchWithTimeout(endpoint, {
      method: method,
      headers: headers,
      body: JSON.stringify(body),
      timeout: 8000
    });
    return [null, res.ok];
  } catch (error: any) {
    console.error(`Error: ${error}`);
    return [null, false];
  }
}
