import React, { useState, useEffect } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import api from '../api/axios';

const RegisterPage = () => {
    const navigate = useNavigate();
    const [groups, setGroups] = useState([]);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState('');
    const [formData, setFormData] = useState({
        first_name: '',
        last_name: '',
        username: '',
        email: '',
        password: '',
        role: 'student',
        group_id: ''
    });

    useEffect(() => {
        const fetchGroups = async () => {
            try {
                const res = await api.groups.getAll();
                setGroups(res.data || []);
            } catch (err) {
                console.error("Ошибка загрузки групп:", err);
            }
        };
        fetchGroups();
    }, []);

    const handleSubmit = async (e) => {
        e.preventDefault();
        setError('');
        setLoading(true);

        try {
            const response = await api.auth.register(formData);
            const { token, user } = response.data.data || response.data;

            if (token?.access_token) {
                localStorage.setItem('access_token', token.access_token);
                if (user) {
                    localStorage.setItem('user_data', JSON.stringify(user));
                }
                navigate('/dashboard', { replace: true });
            } else {
                throw new Error('Токен не получен от сервера');
            }
        } catch (err) {
            setError(err.response?.data?.error || "Ошибка при регистрации");
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className="min-h-screen bg-slate-50 flex items-center justify-center p-6">
            <div className="max-w-md w-full bg-white rounded-[40px] border border-slate-200 p-10 shadow-sm">
                <h2 className="text-3xl font-black text-slate-900 uppercase italic mb-8">
                    Создать <span className="text-brand">аккаунт</span>
                </h2>

                <form onSubmit={handleSubmit} className="space-y-4">
                    {error && (
                        <div className="p-3 bg-red-50 text-red-600 text-xs rounded-xl border border-red-100 italic">
                            {error}
                        </div>
                    )}

                    <div className="grid grid-cols-2 gap-4">
                        <input
                            placeholder="Имя"
                            className="w-full p-4 bg-slate-50 rounded-2xl border-none font-bold text-sm disabled:opacity-50"
                            onChange={e => setFormData({...formData, first_name: e.target.value})}
                            value={formData.first_name}
                            disabled={loading}
                            required
                        />
                        <input
                            placeholder="Фамилия"
                            className="w-full p-4 bg-slate-50 rounded-2xl border-none font-bold text-sm disabled:opacity-50"
                            onChange={e => setFormData({...formData, last_name: e.target.value})}
                            value={formData.last_name}
                            disabled={loading}
                            required
                        />
                    </div>
                    <input
                        placeholder="Username"
                        className="w-full p-4 bg-slate-50 rounded-2xl border-none font-bold text-sm disabled:opacity-50"
                        onChange={e => setFormData({...formData, username: e.target.value})}
                        value={formData.username}
                        disabled={loading}
                        required
                    />
                    <input
                        type="email"
                        placeholder="Email"
                        className="w-full p-4 bg-slate-50 rounded-2xl border-none font-bold text-sm disabled:opacity-50"
                        onChange={e => setFormData({...formData, email: e.target.value})}
                        value={formData.email}
                        disabled={loading}
                        required
                    />
                    <input
                        type="password"
                        placeholder="Пароль"
                        className="w-full p-4 bg-slate-50 rounded-2xl border-none font-bold text-sm disabled:opacity-50"
                        onChange={e => setFormData({...formData, password: e.target.value})}
                        value={formData.password}
                        disabled={loading}
                        required
                    />
                    <select
                        className="w-full p-4 bg-slate-50 rounded-2xl border-none font-bold text-sm text-slate-500 disabled:opacity-50"
                        onChange={e => setFormData({...formData, role: e.target.value})}
                        value={formData.role}
                        disabled={loading}
                    >
                        <option value="student">Студент</option>
                        <option value="teacher">Преподаватель</option>
                    </select>
                    <select
                        value={formData.group_id}
                        onChange={e => setFormData({...formData, group_id: e.target.value})}
                        className="w-full p-4 bg-slate-50 rounded-2xl border-none font-bold text-sm text-slate-500 disabled:opacity-50"
                        disabled={loading}
                    >
                        <option value="">Выберите группу (опционально)</option>
                        {groups.map(g => (
                            <option key={g.id} value={g.id}>{g.name}</option>
                        ))}
                    </select>

                    <button
                        type="submit"
                        disabled={loading}
                        className="w-full bg-brand text-white p-5 rounded-2xl font-black uppercase text-xs tracking-widest shadow-lg shadow-brand/20 mt-4 disabled:bg-slate-300 disabled:cursor-not-allowed"
                    >
                        {loading ? 'Регистрация...' : 'Зарегистрироваться'}
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