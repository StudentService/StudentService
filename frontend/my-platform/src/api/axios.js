import axios from 'axios';

const instance = axios.create({
    baseURL: 'http://localhost:8080/api/v1',
});

// Добавляем токен к КАЖДОМУ запросу
instance.interceptors.request.use((config) => {
    const token = localStorage.getItem('access_token');
    if (token) {
        config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
}, (error) => {
    return Promise.reject(error);
});

// Обработка ответов (ошибки 401 и т.д.)
instance.interceptors.response.use(
    (response) => response,
    (error) => {
        if (error.response && error.response.status === 401) {
            console.error("Токен невалиден. Требуется перезаход.");
            localStorage.removeItem('access_token');
            localStorage.removeItem('user_data');

            // Редирект на логин только если мы не на странице логина/регистрации
            const currentPath = window.location.pathname;
            if (!currentPath.startsWith('/login') && !currentPath.startsWith('/register')) {
                window.location.href = '/login';
            }
        }
        return Promise.reject(error);
    }
);

export default instance;