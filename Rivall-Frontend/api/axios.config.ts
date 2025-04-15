import axios from 'axios';

const HOST = '96.60.10.12'
const PORT = '8080'
const VERSION = 'v1'

// Create an axios instance
const client = axios.create({
    // baseURL: process.env.EXPO_PUBLIC_API_URL + '/api/v1/',
    baseURL: `http://${HOST}:${PORT}/api/${VERSION}/`,
    headers: {
        'Content-Type': 'application/json',
        'Accept': 'application/json'
    },
    timeout: 5000,
});

// client.interceptors.request.use(
//     (config) => {
//         console.log(JSON.stringify(config, null, 2));
//         return config;
//     })

// client.interceptors.response.use(
//     (response) => {
//         console.log('Response:', JSON.stringify(response, null , 3));
//         return response;
//     },
//     (error) => {
//         console.error('Response Error:', JSON.stringify(error, null, 3));
//         return Promise.reject(error);
//     }
// );

export {
    client,
    HOST,
    PORT,
    VERSION
};