import React, { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import api from '../api/axios';

const RegisterPage = () => {
    const [formData, setFormData] = useState({
        email: '',
        first_name: '',
        last_name: '',
        username: '',
        password: '',
        role: 'student'
    });
    const [error, setError] = useState('');
    const navigate = useNavigate();

    const handleSubmit = async (e) => {
        e.preventDefault();
        setError('');

        try {
            const response = await api.post('/auth/register', formData);
            // После регистрации бэкенд возвращает вложенный объект token
            const token = response.data.token?.access_token;

            if (token) {
                localStorage.setItem('access_token', token);
                navigate('/dashboard');
            }
        } catch (err) {
            setError(err.response?.data?.error || 'Ошибка регистрации. Попробуйте другой логин.');
        }
    };

    return (
        <div className="min-h-screen bg-slate-50 flex items-center justify-center p-6 font-sans">
            <div className="w-full max-w-[500px] bg-white rounded-3xl shadow-xl border border-slate-100 p-10">
                <div className="text-center mb-8">
                    <h1 className="text-2xl font-bold text-slate-800 tracking-tight">Создать профиль</h1>
                    <p className="text-slate-500 text-sm mt-1">Присоединяйтесь к кадровому резерву</p>
                </div>

                <form onSubmit={handleSubmit} className="space-y-4">
                    {error && <div className="p-3 bg-red-50 text-red-600 text-xs rounded-xl border border-red-100 italic">{error}</div>}

                    <div className="grid grid-cols-2 gap-4">
                        <input
                            type="text" required placeholder="Имя"
                            className="w-full px-4 py-3 bg-slate-50 border border-slate-200 rounded-2xl outline-none focus:ring-2 focus:ring-blue-500/20 transition-all"
                            onChange={(e) => setFormData({...formData, first_name: e.target.value})}
                        />
                        <input
                            type="text" required placeholder="Фамилия"
                            className="w-full px-4 py-3 bg-slate-50 border border-slate-200 rounded-2xl outline-none focus:ring-2 focus:ring-blue-500/20 transition-all"
                            onChange={(e) => setFormData({...formData, last_name: e.target.value})}
                        />
                    </div>

                    <input
                        type="text" required placeholder="Логин (username)"
                        className="w-full px-4 py-3 bg-slate-50 border border-slate-200 rounded-2xl outline-none focus:ring-2 focus:ring-blue-500/20 transition-all"
                        onChange={(e) => setFormData({...formData, username: e.target.value})}
                    />

                    <input
                        type="email" required placeholder="Email"
                        className="w-full px-4 py-3 bg-slate-50 border border-slate-200 rounded-2xl outline-none focus:ring-2 focus:ring-blue-500/20 transition-all"
                        onChange={(e) => setFormData({...formData, email: e.target.value})}
                    />

                    <input
                        type="password" required placeholder="Пароль"
                        className="w-full px-4 py-3 bg-slate-50 border border-slate-200 rounded-2xl outline-none focus:ring-2 focus:ring-blue-500/20 transition-all"
                        onChange={(e) => setFormData({...formData, password: e.target.value})}
                    />

                    <button type="submit" className="w-full bg-blue-600 hover:bg-blue-700 text-white font-bold py-4 rounded-2xl shadow-lg transition-all active:scale-95">
                        Зарегистрироваться
                    </button>
                </form>

                <p className="mt-8 text-center text-sm text-slate-500">
                    Уже есть аккаунт? <Link to="/login" className="text-blue-600 font-bold hover:underline">Войти</Link>
                </p>
            </div>
        </div>
    );
};

export default RegisterPage;