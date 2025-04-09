import axios from 'axios';

// Create an axios instance
export const client = axios.create({
    baseURL: process.env.EXPO_PUBLIC_API_URL + '/api/v1/',
    timeout: 5000,
});