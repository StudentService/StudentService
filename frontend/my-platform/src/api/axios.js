import axios from 'axios';

const instance = axios.create({
    baseURL: 'http://localhost:8080/api/v1',
});

// Добавляем токен к КАЖДОМУ запросу
instance.interceptors.request.use((config) => {
    const token = localStorage.getItem('access_token');
    if (token) {
        // Убедись, что 'Authorization' написан именно так
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
        }
        return Promise.reject(error);
    }
);
instance.interceptors.response.use(
    (response) => response,
    (error) => {
        if (error.response && error.response.status === 401) {
            // Если сервер сказал "401", значит токен всё.
            localStorage.removeItem('access_token');
            localStorage.removeItem('user'); // почисти всё

            // Жесткий редирект на логин, чтобы юзер не гадал, почему не работает
            window.location.href = '/login';
        }
        return Promise.reject(error);
    }
);

export default instance;