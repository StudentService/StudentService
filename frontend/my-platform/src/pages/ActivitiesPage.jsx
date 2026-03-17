import React, { useState, useEffect } from 'react';
import { api } from '../api';

const ActivitiesPage = () => {
    const [activities, setActivities] = useState([]);
    const [filter, setFilter] = useState('available'); // 'all' или 'available'
    const [loading, setLoading] = useState(true);

    const fetchActivities = async () => {
        try {
            setLoading(true);
            let res;
            if (filter === 'all') {
                // Только для админа
                res = await api.activities.getAll();
            } else {
                res = await api.activities.getAvailable();
            }

            // Преобразуем данные с бэкенда в формат для фронта
            const activitiesData = res.data || [];
            const formattedActivities = activitiesData.map(act => ({
                ...act,
                // Маппинг полей
                category: act.type, // используем type как category
                date: act.start_time || act.created_at,
                participants_count: act.current_participants || 0,
                is_joined: act.is_enrolled || false
            }));

            setActivities(formattedActivities);
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

    const handleLeave = async (id) => {
        try {
            await api.activities.leave(id);
            alert("Запись отменена");
            fetchActivities();
        } catch (err) {
            alert("Не удалось отменить запись");
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
                        onClick={() => setFilter('available')}
                        className={`px-6 py-2 rounded-[14px] text-[10px] font-black uppercase transition-all ${filter === 'available' ? 'bg-white text-[#2D9396] shadow-sm' : 'text-slate-400'}`}
                    >
                        Доступные
                    </button>
                    <button
                        onClick={() => setFilter('all')}
                        className={`px-6 py-2 rounded-[14px] text-[10px] font-black uppercase transition-all ${filter === 'all' ? 'bg-white text-slate-800 shadow-sm' : 'text-slate-400'}`}
                    >
                        Все
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
                                    {act.category || act.type || 'Общее'}
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
                                <span className="opacity-50">📅</span> {act.date ? new Date(act.date).toLocaleDateString() : 'Дата не указана'}
                            </div>
                            <div className="flex items-center gap-3 text-slate-500 font-bold text-[11px]">
                                <span className="opacity-50">📍</span> {act.location || 'Главный корпус'}
                            </div>
                            <div className="flex items-center gap-3 text-slate-500 font-bold text-[11px]">
                                <span className="opacity-50">👥</span> Мест: {act.participants_count || 0} / {act.max_participants || '∞'}
                            </div>
                        </div>

                        {/* Кнопка записи */}
                        <div className="p-6">
                            {act.is_joined ? (
                                <button
                                    onClick={() => handleLeave(act.id)}
                                    className="w-full py-4 rounded-2xl font-black uppercase text-[10px] tracking-widest bg-red-50 text-red-500 border border-red-100 hover:bg-red-100 transition-all"
                                >
                                    Отменить запись
                                </button>
                            ) : (
                                <button
                                    onClick={() => handleJoin(act.id)}
                                    disabled={act.participants_count >= act.max_participants}
                                    className={`w-full py-4 rounded-2xl font-black uppercase text-[10px] tracking-widest transition-all ${
                                        act.participants_count >= act.max_participants
                                            ? 'bg-slate-50 text-slate-400 cursor-not-allowed'
                                            : 'bg-[#2D9396] text-white shadow-lg shadow-[#2D9396]/20 hover:scale-[1.02] active:scale-95'
                                    }`}
                                >
                                    {act.participants_count >= act.max_participants ? 'Мест нет' : 'Записаться'}
                                </button>
                            )}
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