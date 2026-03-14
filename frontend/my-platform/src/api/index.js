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
    questionnaires: {
        getMy: () => apiAxios.get('/questionnaire/my'),
        getTemplate: () => apiAxios.get('/questionnaire/template'),
        submit: (answers) => apiAxios.post('/questionnaire/submit', { answers }),
    },
    challenges: {
        getMy: () => apiAxios.get('/challenges/my'),
        getAllAvailable: () => apiAxios.get('/activities/available'),
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
    }
};