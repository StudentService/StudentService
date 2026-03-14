import React, { useState, useEffect } from 'react';
import { api } from '../api';

const ChallengesPage = () => {
    const [challenges, setChallenges] = useState([]);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        const fetchChallenges = async () => {
            try {
                const res = await api.challenges.getMy();
                setChallenges(res.data.data || res.data || []);
            } catch (err) {
                console.error("Ошибка загрузки вызовов:", err);
            } finally {
                setLoading(false);
            }
        };
        fetchChallenges();
    }, []);

    if (loading) return <div className="p-10 font-black text-[#2D9396] animate-pulse">СИНХРОНИЗАЦИЯ ВЫЗОВОВ...</div>;

    return (
        <div className="p-8 space-y-8 animate-in fade-in duration-500">
            <div className="flex justify-between items-center">
                <h1 className="text-3xl font-black text-slate-800 uppercase italic tracking-tighter">Мои вызовы</h1>
                <button className="bg-slate-900 text-white px-6 py-3 rounded-2xl font-black text-[10px] uppercase tracking-widest">
                    + Новый вызов
                </button>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                {challenges.map((ch) => (
                    <div key={ch.id} className="bg-white rounded-[32px] border border-slate-100 p-8 shadow-sm hover:border-[#2D9396] transition-all cursor-pointer group">
                        <div className="flex justify-between items-start mb-6">
                            <span className={`text-[10px] font-black uppercase px-3 py-1 rounded-lg ${ch.status === 'active' ? 'bg-green-50 text-green-500' : 'bg-slate-100 text-slate-400'}`}>
                                {ch.status}
                            </span>
                            <span className="text-[#2D9396] font-black text-xl">{ch.progress}%</span>
                        </div>

                        <h3 className="text-xl font-bold text-slate-800 mb-2 group-hover:text-[#2D9396] transition-colors">{ch.title}</h3>
                        <p className="text-slate-500 text-sm line-clamp-2 mb-6">{ch.description || "Описание отсутствует"}</p>

                        <div className="space-y-4">
                            <div className="h-2 bg-slate-50 rounded-full overflow-hidden">
                                <div className="h-full bg-[#2D9396]" style={{ width: `${ch.progress}%` }}></div>
                            </div>
                            <div className="flex justify-between text-[10px] font-black text-slate-400 uppercase tracking-tighter">
                                <span>Начало: {new Date(ch.start_date).toLocaleDateString()}</span>
                                <span>Конец: {new Date(ch.end_date).toLocaleDateString()}</span>
                            </div>
                        </div>
                    </div>
                ))}
            </div>

            {challenges.length === 0 && (
                <div className="text-center py-20 bg-slate-50 rounded-[40px] border-2 border-dashed border-slate-100">
                    <p className="text-slate-400 font-bold uppercase text-xs">Список вызовов пуст</p>
                </div>
            )}
        </div>
    );
};

export default ChallengesPage;