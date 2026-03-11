import axios from 'axios';

const api = axios.create({
    baseURL: 'http://localhost:8080/api/v1',
});

// Логируем ошибки в консоль для удобства отладки
api.interceptors.response.use(
    (response) => response,
    (error) => {
        console.error('Ошибка API:', error.response?.data || error.message);
        return Promise.reject(error);
    }
);

api.interceptors.request.use((config) => {
    const token = localStorage.getItem('access_token');
    if (token) {
        config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
});

export default api;