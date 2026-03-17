import React, { useState, useEffect } from 'react';
import { api } from '../api';

const ProfilePage = () => {
    const [user, setUser] = useState(null);
    const [loading, setLoading] = useState(true);
    const [editing, setEditing] = useState(false);
    const [formData, setFormData] = useState({});

    useEffect(() => {
        const fetchProfile = async () => {
            try {
                const res = await api.users.getMe();
                const userData = res.data.data || res.data;
                setUser(userData);
                setFormData({
                    first_name: userData.first_name || '',
                    last_name: userData.last_name || '',
                    username: userData.username || '',
                });
            } catch (err) {
                console.error("Profile error", err);
            } finally {
                setLoading(false);
            }
        };
        fetchProfile();
    }, []);

    const handleInputChange = (e) => {
        const { name, value } = e.target;
        setFormData(prev => ({ ...prev, [name]: value }));
    };

    const handleSubmit = async (e) => {
        e.preventDefault();
        try {
            const res = await api.users.updateMe(formData);
            const updatedUser = res.data.data || res.data;
            setUser(updatedUser);
            setEditing(false);
        } catch (err) {
            console.error("Update error", err);
            alert("Ошибка при обновлении профиля");
        }
    };

    if (loading) return (
        <div className="p-10 font-black text-[#2D9396] animate-pulse uppercase text-center">
            Загрузка профиля...
        </div>
    );

    // Определяем цвет и текст для роли
    const getRoleInfo = (role) => {
        switch(role) {
            case 'admin':
                return { color: 'bg-purple-100 text-purple-600', text: 'Администратор' };
            case 'teacher':
                return { color: 'bg-blue-100 text-blue-600', text: 'Преподаватель' };
            case 'holder':
                return { color: 'bg-green-100 text-green-600', text: 'Держатель' };
            case 'student':
                return { color: 'bg-[#2D9396]/10 text-[#2D9396]', text: 'Студент' };
            default:
                return { color: 'bg-slate-100 text-slate-600', text: role };
        }
    };

    const roleInfo = getRoleInfo(user?.role);

    return (
        <div className="p-4 md:p-8 max-w-5xl mx-auto animate-in fade-in duration-500">
            {/* Заголовок */}
            <header className="mb-8">
                <h1 className="text-4xl font-black text-slate-800 uppercase italic tracking-tighter">
                    Личный кабинет
                </h1>
                <p className="text-slate-400 font-bold text-[10px] uppercase tracking-[0.2em] mt-1">
                    Управление профилем
                </p>
            </header>

            <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
                {/* Левая колонка - аватар и основная информация */}
                <div className="lg:col-span-1">
                    <div className="bg-white rounded-[32px] border border-slate-100 shadow-sm overflow-hidden sticky top-8">
                        {/* Шапка с градиентом */}
                        <div className="h-32 bg-gradient-to-r from-[#2D9396] to-[#40b8bb] relative">
                            <div className="absolute -bottom-12 left-8">
                                <div className="w-24 h-24 bg-white rounded-[24px] shadow-lg flex items-center justify-center text-4xl font-black text-[#2D9396] border-4 border-white">
                                    {user?.first_name?.[0]}{user?.last_name?.[0]}
                                </div>
                            </div>
                        </div>

                        {/* Информация под аватаром */}
                        <div className="pt-16 p-8">
                            <div className="mb-6">
                                <h2 className="text-2xl font-bold text-slate-800">
                                    {user?.first_name} {user?.last_name}
                                </h2>
                                <p className="text-slate-400 text-sm mt-1">@{user?.username}</p>
                            </div>

                            {/* Роль */}
                            <div className="mb-8">
                                <span className={`inline-block px-4 py-2 rounded-xl text-[10px] font-black uppercase ${roleInfo.color}`}>
                                    {roleInfo.text}
                                </span>
                            </div>

                            {/* Статистика */}
                            <div className="grid grid-cols-2 gap-4 pt-6 border-t border-slate-100">
                                <div>
                                    <p className="text-[10px] font-black text-slate-400 uppercase">Создан</p>
                                    <p className="font-bold text-slate-700 text-sm">
                                        {new Date(user?.created_at).toLocaleDateString()}
                                    </p>
                                </div>
                                <div>
                                    <p className="text-[10px] font-black text-slate-400 uppercase">Обновлён</p>
                                    <p className="font-bold text-slate-700 text-sm">
                                        {new Date(user?.updated_at).toLocaleDateString()}
                                    </p>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>

                {/* Правая колонка - детальная информация и редактирование */}
                <div className="lg:col-span-2">
                    <div className="bg-white rounded-[32px] border border-slate-100 shadow-sm overflow-hidden">
                        {/* Заголовок с кнопкой редактирования */}
                        <div className="p-8 border-b border-slate-100 flex justify-between items-center">
                            <h3 className="text-sm font-black text-slate-400 uppercase tracking-widest">
                                Детальная информация
                            </h3>
                            {!editing && (
                                <button
                                    onClick={() => setEditing(true)}
                                    className="px-6 py-3 bg-slate-50 hover:bg-slate-100 rounded-2xl text-[10px] font-black uppercase transition-all"
                                >
                                    ✎ Редактировать
                                </button>
                            )}
                        </div>

                        {editing ? (
                            // Форма редактирования
                            <form onSubmit={handleSubmit} className="p-8 space-y-6">
                                <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                                    <div>
                                        <label className="text-[10px] font-black text-slate-400 uppercase mb-2 block ml-2">
                                            Имя
                                        </label>
                                        <input
                                            type="text"
                                            name="first_name"
                                            value={formData.first_name}
                                            onChange={handleInputChange}
                                            className="w-full bg-slate-50 border-none rounded-2xl p-4 focus:ring-2 focus:ring-[#2D9396] transition-all"
                                            required
                                        />
                                    </div>
                                    <div>
                                        <label className="text-[10px] font-black text-slate-400 uppercase mb-2 block ml-2">
                                            Фамилия
                                        </label>
                                        <input
                                            type="text"
                                            name="last_name"
                                            value={formData.last_name}
                                            onChange={handleInputChange}
                                            className="w-full bg-slate-50 border-none rounded-2xl p-4 focus:ring-2 focus:ring-[#2D9396] transition-all"
                                            required
                                        />
                                    </div>
                                </div>

                                <div>
                                    <label className="text-[10px] font-black text-slate-400 uppercase mb-2 block ml-2">
                                        Username
                                    </label>
                                    <input
                                        type="text"
                                        name="username"
                                        value={formData.username}
                                        onChange={handleInputChange}
                                        className="w-full bg-slate-50 border-none rounded-2xl p-4 focus:ring-2 focus:ring-[#2D9396] transition-all"
                                        required
                                    />
                                </div>

                                <div className="bg-slate-50 p-6 rounded-3xl">
                                    <p className="text-[10px] font-black text-slate-400 uppercase mb-1">Email</p>
                                    <p className="font-bold text-slate-700">{user?.email}</p>
                                    <p className="text-[8px] text-slate-400 mt-2 italic">Email нельзя изменить</p>
                                </div>

                                <div className="flex gap-4 pt-4">
                                    <button
                                        type="submit"
                                        className="flex-1 bg-[#2D9396] text-white py-4 rounded-2xl font-black uppercase text-[10px] tracking-widest shadow-lg shadow-[#2D9396]/20 hover:scale-[1.01] transition-all"
                                    >
                                        Сохранить
                                    </button>
                                    <button
                                        type="button"
                                        onClick={() => setEditing(false)}
                                        className="px-8 py-4 bg-slate-50 hover:bg-slate-100 rounded-2xl font-black uppercase text-[10px] tracking-widest transition-all"
                                    >
                                        Отмена
                                    </button>
                                </div>
                            </form>
                        ) : (
                            // Просмотр информации
                            <div className="p-8 space-y-6">
                                <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                                    <div className="bg-slate-50 p-6 rounded-3xl">
                                        <p className="text-[10px] font-black text-slate-400 uppercase mb-1">Имя</p>
                                        <p className="font-bold text-slate-700 text-lg">{user?.first_name}</p>
                                    </div>
                                    <div className="bg-slate-50 p-6 rounded-3xl">
                                        <p className="text-[10px] font-black text-slate-400 uppercase mb-1">Фамилия</p>
                                        <p className="font-bold text-slate-700 text-lg">{user?.last_name}</p>
                                    </div>
                                </div>

                                <div className="bg-slate-50 p-6 rounded-3xl">
                                    <p className="text-[10px] font-black text-slate-400 uppercase mb-1">Username</p>
                                    <p className="font-bold text-slate-700">@{user?.username}</p>
                                </div>

                                <div className="bg-slate-50 p-6 rounded-3xl">
                                    <p className="text-[10px] font-black text-slate-400 uppercase mb-1">Email</p>
                                    <p className="font-bold text-slate-700">{user?.email}</p>
                                </div>

                                <div className="bg-slate-50 p-6 rounded-3xl">
                                    <p className="text-[10px] font-black text-slate-400 uppercase mb-1">Группа</p>
                                    <p className="font-bold text-slate-700">{user?.group_id || 'Не указана'}</p>
                                </div>
                            </div>
                        )}
                    </div>
                </div>
            </div>
        </div>
    );
};

export default ProfilePage;