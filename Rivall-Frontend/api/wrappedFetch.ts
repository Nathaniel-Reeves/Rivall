
export async function wrappedFetch(endpoint: string, method: string, headers: Record<string, string>, body: Object) : [Promise<Object>, Number] {

    // Validate Endpoint
    if (typeof endpoint !== 'string') {
      throw new Error(`Invalid endpoint: ${endpoint}`);
    }
  
    // Validate Method (GET, POST, PUT, DELETE)
    if (!['GET', 'POST', 'PUT', 'DELETE'].includes(method)) {
      throw new Error(`Invalid method: ${method}`);
    }
  
    // Validate Headers
    if (typeof headers !== 'object') {
      throw new Error(`Invalid headers: ${headers}`);
    }
    const h: Record<string, string> = {
      'Content-Type': 'application/json',
      Accept: 'application/json',
      ...headers,
    };

    console.debug(`fetching: ${method} ${endpoint}`);
    console.debug('headers:', JSON.stringify(h, null, 2));
    let frame: {
      method: string;
      headers: Record<string, string>;
      body?: string;
    } = {
      method: method,
      headers: h,
    }
  
    // Validate Body
    if (typeof body !== 'object' || Array.isArray(body)) {
      throw new Error(`Invalid body: ${body}`);
    }
    const b = body ? JSON.stringify(body) : '';
    if (b !== '' && method !== 'GET') {
      console.debug('body:', JSON.stringify(b, null, 2));
      frame = {
        ...frame,
        body: b,
      }
    }

    return fetch(endpoint, frame)
    .then(response => {
      if (!response.ok) {
        console.warn('Network response was not ok');
        console.warn('response status:', response.status);
        return [{ error: response.statusText, status: response.status }, response.status];
      }
      return [response.json().then(data => {
          console.debug(`response data: ${JSON.stringify(data, null, 2)}`);
          return data;
      }), response.status];
    })
    .catch((error) => {
        console.error('Error:', error);
        return [{ error: error.message }, 500];
    });
  }