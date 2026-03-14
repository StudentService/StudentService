import axios from 'axios';

const instance = axios.create({
    baseURL: 'http://localhost:8080/api/v1', // Убедись, что адрес бэкенда верный
});

// Интерцептор для добавления токена
instance.interceptors.request.use((config) => {
    const token = localStorage.getItem('access_token');
    if (token) {
        config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
});

export default instance;