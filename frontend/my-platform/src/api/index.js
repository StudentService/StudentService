import apiAxios from './axios';

export const api = {
    auth: {
        login: (credentials) => apiAxios.post('/auth/login', credentials),
        register: (data) => apiAxios.post('/auth/register', data),
    },
    users: {
        getMe: () => apiAxios.get('/users/me'),
        updateMe: (data) => apiAxios.patch('/users/me', data), // ДОБАВЛЕНО
    },
    dashboard: {
        getSummary: () => apiAxios.get('/dashboard'),
    },
    challenges: {
        getMy: () => apiAxios.get('/challenges/my'),
        create: (data) => apiAxios.post('/challenges', data),
    },
    calendar: {
        getEvents: (from, to) => apiAxios.get('/calendar/events/my', {
            params: { from, to }
        }),
    },
    grades: {
        getSummary: () => apiAxios.get('/grades/my/summary'),
        getAll: () => apiAxios.get('/grades/my'),
        getByCourse: (courseId) => apiAxios.get(`/grades/my/courses/${courseId}`),
    },
    activities: {
        getAll: () => apiAxios.get('/activities'),
        getAvailable: () => apiAxios.get('/activities/available'),
        getMy: () => apiAxios.get('/activities/my'),
        join: (activityId) => apiAxios.post(`/activities/${activityId}/enroll`),
        leave: (activityId) => apiAxios.delete(`/activities/${activityId}/enroll`),
    },
    questionnaires: {
        getMy: () => apiAxios.get('/questionnaire/my'), // ДОБАВЛЕНО
        getTemplate: () => apiAxios.get('/questionnaire/template'),
        submit: (data) => apiAxios.post('/questionnaire/submit', { answers: data }),
        saveDraft: (data) => apiAxios.post('/questionnaire/draft', data), // ДОБАВЛЕНО
    }
};