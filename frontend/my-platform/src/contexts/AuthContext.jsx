import React, { createContext, useContext, useState, useEffect } from 'react';
import { api } from '../api';

const AuthContext = createContext(null);

export const useAuth = () => {
    const context = useContext(AuthContext);
    if (!context) {
        throw new Error('useAuth must be used within an AuthProvider');
    }
    return context;
};

export const AuthProvider = ({ children }) => {
    const [user, setUser] = useState(null);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        loadUser();
    }, []);

    const loadUser = async () => {
        const token = localStorage.getItem('access_token');
        const storedUser = localStorage.getItem('user_data');

        if (!token) {
            setLoading(false);
            return;
        }

        try {
            // Если есть сохранённые данные пользователя, используем их
            if (storedUser) {
                setUser(JSON.parse(storedUser));
            } else {
                // Загружаем данные с сервера
                const response = await api.users.getMe();
                const userData = response.data;
                setUser(userData);
                localStorage.setItem('user_data', JSON.stringify(userData));
            }
        } catch (error) {
            console.error('Failed to load user:', error);
            localStorage.removeItem('access_token');
            localStorage.removeItem('user_data');
            setUser(null);
        } finally {
            setLoading(false);
        }
    };

    const login = async (credentials) => {
        const response = await api.auth.login(credentials);
        const { access_token, refresh_token } = response.data;

        localStorage.setItem('access_token', access_token);
        localStorage.setItem('refresh_token', refresh_token);

        // Загружаем данные пользователя
        const userResponse = await api.users.getMe();
        const userData = userResponse.data;
        setUser(userData);
        localStorage.setItem('user_data', JSON.stringify(userData));

        return userData;
    };

    const register = async (data) => {
        const response = await api.auth.register(data);
        const { access_token, refresh_token, user } = response.data;

        localStorage.setItem('access_token', access_token);
        localStorage.setItem('refresh_token', refresh_token);
        localStorage.setItem('user_data', JSON.stringify(user));

        setUser(user);
        return user;
    };

    const logout = () => {
        localStorage.removeItem('access_token');
        localStorage.removeItem('refresh_token');
        localStorage.removeItem('user_data');
        setUser(null);
    };

    const updateUser = (updatedData) => {
        const updatedUser = { ...user, ...updatedData };
        setUser(updatedUser);
        localStorage.setItem('user_data', JSON.stringify(updatedUser));
    };

    // Проверки ролей
    const isStudent = () => user?.role === 'student';
    const isTeacher = () => user?.role === 'teacher';
    const isHolder = () => user?.role === 'holder';
    const isCandidate = () => user?.role === 'candidate';
    const isAdmin = () => user?.role === 'admin';

    const hasRole = (roles) => {
        if (!user) return false;
        if (Array.isArray(roles)) {
            return roles.includes(user.role);
        }
        return user.role === roles;
    };

    const value = {
        user,
        loading,
        login,
        register,
        logout,
        updateUser,
        isStudent,
        isTeacher,
        isHolder,
        isCandidate,
        isAdmin,
        hasRole,
        isAuthenticated: !!user,
    };

    return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
};