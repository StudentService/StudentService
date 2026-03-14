import apiAxios from './axios';

export const api = {
    auth: {
        login: (credentials) => apiAxios.post('/auth/login', credentials),
        register: (data) => apiAxios.post('/auth/register', data),
    },
    users: {
        getMe: () => apiAxios.get('/users/me'),
    },
    dashboard: {
        getSummary: () => apiAxios.get('/dashboard'),
    },
    challenges: {
        getMy: () => apiAxios.get('/challenges'), // Метод GetMyChallenges
        create: (data) => apiAxios.post('/challenges', data), // Метод CreateChallenge
    },
    events: {
        getStats: () => apiAxios.get('/dashboard'), // В твоем эндпоинте dashboard есть upcoming_events_count
        getAll: () => apiAxios.get('/calendar/events'), // Предполагаемый эндпоинт для списка событий
    },
    grades: {
        getSummary: () => apiAxios.get('/grades/my/summary'),
        getAll: () => apiAxios.get('/grades/my'),
        getByCourse: (courseId) => apiAxios.get(`/grades/my/courses/${courseId}`),
    },
    activities: {
        getAll: () => apiAxios.get('/activities'),
        getAvailable: () => apiAxios.get('/activities/available'),
        join: (activityId) => apiAxios.post(`/activities/${activityId}/join`),
        leave: (activityId) => apiAxios.delete(`/activities/${activityId}/leave`),
    },
    questionnaires: {
        getTemplate: () => apiAxios.get('/questionnaire/template'),
        submit: (data) => apiAxios.post('/questionnaire/submit', {
            answers: data // бэкенд обычно ждет массив объектов в поле answers
        }),
    }
};