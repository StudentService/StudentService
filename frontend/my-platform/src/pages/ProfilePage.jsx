import React, { useState, useEffect } from 'react';
import { api } from '../api';

const ProfilePage = () => {
    const [user, setUser] = useState(null);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        const fetchProfile = async () => {
            try {
                // Прямой вызов через наш объект api
                const res = await api.users.getMe();
                setUser(res.data.data || res.data);
            } catch (err) {
                console.error("Profile error", err);
            } finally {
                setLoading(false);
            }
        };
        fetchProfile();
    }, []);

    if (loading) return <div className="p-10 font-black text-[#2D9396] animate-pulse uppercase">Загрузка профиля...</div>;

    return (
        <div className="max-w-4xl p-8 animate-in fade-in duration-500">
            <h1 className="text-3xl font-black text-slate-800 uppercase italic mb-8 italic">Личный кабинет</h1>
            <div className="bg-white rounded-[32px] border border-slate-100 shadow-sm overflow-hidden">
                <div className="bg-slate-50 p-8 flex items-center gap-6">
                    <div className="w-20 h-20 bg-[#2D9396] rounded-[24px] flex items-center justify-center text-white text-2xl font-black italic">
                        {user?.first_name?.[0]}
                    </div>
                    <div>
                        <h2 className="text-2xl font-bold text-slate-800">{user?.first_name} {user?.last_name}</h2>
                        <p className="text-slate-400 font-bold text-xs uppercase tracking-widest italic">{user?.role}</p>
                    </div>
                </div>
                <div className="p-8 grid grid-cols-1 md:grid-cols-2 gap-6">
                    <div className="p-4 bg-slate-50 rounded-2xl">
                        <p className="text-[10px] font-black text-slate-400 uppercase mb-1">Email</p>
                        <p className="font-bold text-slate-800">{user?.email}</p>
                    </div>
                    <div className="p-4 bg-slate-50 rounded-2xl">
                        <p className="text-[10px] font-black text-slate-400 uppercase mb-1">Группа</p>
                        <p className="font-bold text-slate-800">{user?.group_id || 'Не указана'}</p>
                    </div>
                </div>
            </div>
        </div>
    );
};

export default ProfilePage;