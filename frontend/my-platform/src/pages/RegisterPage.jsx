import React, { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { api } from '../api';

const RegisterPage = () => {
    const navigate = useNavigate();
    const [formData, setFormData] = useState({
        first_name: '',
        last_name: '',
        username: '',
        email: '',
        password: '',
        role: 'student' // Значение по умолчанию
    });

    const handleSubmit = async (e) => {
        e.preventDefault();
        try {
            const response = await api.auth.register(formData);
            // Согласно твоему Swagger: ответ содержит { token: { access_token... }, user: { ... } }
            const { token, user } = response.data.data || response.data;

            if (token?.access_token) {
                localStorage.setItem('access_token', token.access_token);
                localStorage.setItem('user_data', JSON.stringify(user));
                navigate('/dashboard');
            }
        } catch (err) {
            alert(err.response?.data?.error || "Ошибка при регистрации");
        }
    };

    return (
        <div className="min-h-screen bg-slate-50 flex items-center justify-center p-6">
            <div className="max-w-md w-full bg-white rounded-[40px] border border-slate-200 p-10 shadow-sm">
                <h2 className="text-3xl font-black text-slate-900 uppercase italic mb-8">
                    Создать <span className="text-brand">аккаунт</span>
                </h2>

                <form onSubmit={handleSubmit} className="space-y-4">
                    <div className="grid grid-cols-2 gap-4">
                        <input
                            placeholder="Имя"
                            className="p-4 bg-slate-50 rounded-2xl border-none font-bold text-sm"
                            onChange={e => setFormData({...formData, first_name: e.target.value})}
                            required
                        />
                        <input
                            placeholder="Фамилия"
                            className="p-4 bg-slate-50 rounded-2xl border-none font-bold text-sm"
                            onChange={e => setFormData({...formData, last_name: e.target.value})}
                            required
                        />
                    </div>
                    <input
                        placeholder="Username"
                        className="w-full p-4 bg-slate-50 rounded-2xl border-none font-bold text-sm"
                        onChange={e => setFormData({...formData, username: e.target.value})}
                        required
                    />
                    <input
                        type="email"
                        placeholder="Email"
                        className="w-full p-4 bg-slate-50 rounded-2xl border-none font-bold text-sm"
                        onChange={e => setFormData({...formData, email: e.target.value})}
                        required
                    />
                    <input
                        type="password"
                        placeholder="Пароль"
                        className="w-full p-4 bg-slate-50 rounded-2xl border-none font-bold text-sm"
                        onChange={e => setFormData({...formData, password: e.target.value})}
                        required
                    />
                    <select
                        className="w-full p-4 bg-slate-50 rounded-2xl border-none font-bold text-sm text-slate-500"
                        onChange={e => setFormData({...formData, role: e.target.value})}
                    >
                        <option value="student">Студент</option>
                        <option value="teacher">Преподаватель</option>
                    </select>

                    <button type="submit" className="w-full bg-brand text-white p-5 rounded-2xl font-black uppercase text-xs tracking-widest shadow-lg shadow-brand/20 mt-4">
                        Зарегистрироваться
                    </button>
                </form>

                <p className="text-center mt-6 text-slate-400 font-bold text-xs uppercase">
                    Уже есть аккаунт? <Link to="/login" className="text-brand">Войти</Link>
                </p>
            </div>
        </div>
    );
};

export default RegisterPage;