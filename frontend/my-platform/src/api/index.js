import apiAxios from './axios';

export const api = {
    auth: {
        login: (credentials) => apiAxios.post('/auth/login', credentials),
        register: (data) => apiAxios.post('/auth/register', data),
    },
    users: {
        getMe: () => apiAxios.get('/users/me'),
        updateMe: (data) => apiAxios.patch('/users/me', data),
    },
    dashboard: {
        getStudent: () => apiAxios.get('/dashboard'),
        getTeacher: () => apiAxios.get('/teacher/dashboard'),
    },
    teacher: {
        getDashboard: () => apiAxios.get('/teacher/dashboard'),
        getGroups: () => apiAxios.get('/teacher/groups'),
        getGroupStudents: (groupId) => apiAxios.get(`/teacher/groups/${groupId}/students`),
        getGroupGrades: (groupId) => apiAxios.get(`/teacher/groups/${groupId}/grades`),
        getStudentProfile: (studentId) => apiAxios.get(`/teacher/students/${studentId}`),
        getStudentGrades: (studentId) => apiAxios.get(`/teacher/students/${studentId}/grades`),
        getActivities: () => apiAxios.get('/teacher/activities'),
        markAttendance: (activityId, data) => apiAxios.post(`/teacher/activities/${activityId}/attendance`, data),
        importGrades: (data) => apiAxios.post('/teacher/grades/import', data),
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
        getMy: () => apiAxios.get('/questionnaire/my'),
        getTemplate: () => apiAxios.get('/questionnaire/template'),
        submit: (data) => apiAxios.post('/questionnaire/submit', { answers: data }),
        saveDraft: (data) => apiAxios.post('/questionnaire/draft', data),
    },
    groups: {
        getAll: () => apiAxios.get('/groups'),
        getById: (id) => apiAxios.get(`/groups/${id}`),
        create: (data) => apiAxios.post('/groups', data),
        update: (id, data) => apiAxios.patch(`/groups/${id}`, data),
        delete: (id) => apiAxios.delete(`/groups/${id}`),
    },
    semesters: {
        getAll: () => apiAxios.get('/semesters'),
        getActive: () => apiAxios.get('/semesters/active'),
        getById: (id) => apiAxios.get(`/semesters/${id}`),
        create: (data) => apiAxios.post('/semesters', data),
        update: (id, data) => apiAxios.patch(`/semesters/${id}`, data),
        delete: (id) => apiAxios.delete(`/semesters/${id}`),
    }
};