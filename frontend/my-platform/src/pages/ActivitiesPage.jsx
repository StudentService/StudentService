import React, { useState, useEffect } from 'react';
import { api } from '../api';

const ActivitiesPage = () => {
    const [activities, setActivities] = useState([]);
    const [filter, setFilter] = useState('all'); // 'all' или 'available'
    const [loading, setLoading] = useState(true);

    const fetchActivities = async () => {
        try {
            setLoading(true);
            const res = filter === 'all'
                ? await api.activities.getAll()
                : await api.activities.getAvailable();
            setActivities(res.data.data || res.data || []);
        } catch (err) {
            console.error("Ошибка загрузки активностей:", err);
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchActivities();
    }, [filter]);

    const handleJoin = async (id) => {
        try {
            await api.activities.join(id);
            alert("Вы успешно записаны!");
            fetchActivities(); // Обновляем список
        } catch (err) {
            alert("Не удалось записаться на активность");
        }
    };

    if (loading) return <div className="p-10 font-black text-[#2D9396] animate-pulse uppercase italic">Поиск событий...</div>;

    return (
        <div className="p-4 md:p-8 space-y-8 animate-in fade-in duration-500">
            <header className="flex flex-col md:flex-row justify-between items-start md:items-center gap-6">
                <div>
                    <h1 className="text-4xl font-black text-slate-800 uppercase italic tracking-tighter">Активности</h1>
                    <p className="text-slate-400 font-bold text-[10px] uppercase tracking-[0.2em] mt-1">Мероприятия, воркшопы и спорт</p>
                </div>

                {/* Переключатель фильтров */}
                <div className="flex bg-slate-100 p-1.5 rounded-[20px]">
                    <button
                        onClick={() => setFilter('all')}
                        className={`px-6 py-2 rounded-[14px] text-[10px] font-black uppercase transition-all ${filter === 'all' ? 'bg-white text-slate-800 shadow-sm' : 'text-slate-400'}`}
                    >
                        Все
                    </button>
                    <button
                        onClick={() => setFilter('available')}
                        className={`px-6 py-2 rounded-[14px] text-[10px] font-black uppercase transition-all ${filter === 'available' ? 'bg-white text-[#2D9396] shadow-sm' : 'text-slate-400'}`}
                    >
                        Доступные
                    </button>
                </div>
            </header>

            <div className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-6">
                {activities.map((act) => (
                    <div key={act.id} className="bg-white rounded-[32px] border border-slate-100 overflow-hidden shadow-sm hover:shadow-md transition-all group">
                        {/* Плашка с типом */}
                        <div className="p-8 pb-0">
                            <div className="flex justify-between items-start mb-6">
                                <span className="bg-slate-50 text-slate-400 text-[9px] font-black uppercase px-3 py-1.5 rounded-lg tracking-widest">
                                    {act.category || 'Общее'}
                                </span>
                                <div className="flex flex-col items-end">
                                    <span className="text-[#FF7A50] font-black text-lg">+{act.points || 0}</span>
                                    <span className="text-[8px] font-black text-slate-300 uppercase">баллов</span>
                                </div>
                            </div>

                            <h3 className="text-xl font-bold text-slate-800 leading-tight mb-2 group-hover:text-[#2D9396] transition-colors">
                                {act.title}
                            </h3>
                            <p className="text-slate-400 text-xs line-clamp-2 font-medium mb-6">
                                {act.description}
                            </p>
                        </div>

                        {/* Детали */}
                        <div className="px-8 py-6 bg-slate-50/50 space-y-3">
                            <div className="flex items-center gap-3 text-slate-500 font-bold text-[11px]">
                                <span className="opacity-50">📅</span> {new Date(act.date).toLocaleDateString()}
                            </div>
                            <div className="flex items-center gap-3 text-slate-500 font-bold text-[11px]">
                                <span className="opacity-50">📍</span> {act.location || 'Главный корпус'}
                            </div>
                            <div className="flex items-center gap-3 text-slate-500 font-bold text-[11px]">
                                <span className="opacity-50">👥</span> Мест: {act.participants_count} / {act.max_participants}
                            </div>
                        </div>

                        {/* Кнопка записи */}
                        <div className="p-6">
                            <button
                                onClick={() => handleJoin(act.id)}
                                disabled={act.is_joined || act.participants_count >= act.max_participants}
                                className={`w-full py-4 rounded-2xl font-black uppercase text-[10px] tracking-widest transition-all ${
                                    act.is_joined
                                        ? 'bg-green-50 text-green-500 border border-green-100'
                                        : 'bg-[#2D9396] text-white shadow-lg shadow-[#2D9396]/20 hover:scale-[1.02] active:scale-95'
                                }`}
                            >
                                {act.is_joined ? 'Вы участвуете' : 'Записаться'}
                            </button>
                        </div>
                    </div>
                ))}
            </div>

            {activities.length === 0 && (
                <div className="text-center py-24 bg-slate-50 rounded-[40px] border-2 border-dashed border-slate-100">
                    <p className="text-slate-300 font-black uppercase text-xs tracking-[0.2em]">На данный момент событий нет</p>
                </div>
            )}
        </div>
    );
};

export default ActivitiesPage;