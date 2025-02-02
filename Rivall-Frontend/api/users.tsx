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
