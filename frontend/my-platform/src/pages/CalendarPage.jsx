import React, { useState, useEffect } from 'react';
import { api } from '../api';

const CalendarPage = () => {
    const [currentDate, setCurrentDate] = useState(new Date());
    const [events, setEvents] = useState([]);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        const fetchEvents = async () => {
            try {
                // Здесь мы можем тянуть данные из API
                // Для примера используем dashboard или специальный эндпоинт событий
                const res = await api.dashboard.getSummary();
                // Если API отдает массив событий, сохраняем его
                setEvents(res.data.upcoming_events || []);
            } catch (err) {
                console.error("Ошибка загрузки календаря:", err);
            } finally {
                setLoading(false);
            }
        };
        fetchEvents();
    }, []);

    const daysInMonth = (year, month) => new Date(year, month + 1, 0).getDate();
    const firstDayOfMonth = new Date(currentDate.getFullYear(), currentDate.getMonth(), 1).getDay();

    const days = [];
    // Пустые ячейки для начала месяца (смещение)
    for (let i = 0; i < (firstDayOfMonth === 0 ? 6 : firstDayOfMonth - 1); i++) {
        days.push(null);
    }
    // Заполняем числами
    for (let i = 1; i <= daysInMonth(currentDate.getFullYear(), currentDate.getMonth()); i++) {
        days.push(i);
    }

    const monthNames = ["Январь", "Февраль", "Март", "Апрель", "Май", "Июнь", "Июль", "Август", "Сентябрь", "Октябрь", "Ноябрь", "Декабрь"];

    if (loading) return <div className="p-10 font-black text-[#2D9396] animate-pulse uppercase">Загрузка расписания...</div>;

    return (
        <div className="p-4 md:p-8 space-y-8 animate-in fade-in duration-500">
            <header className="flex flex-col md:flex-row justify-between items-start md:items-center gap-4">
                <div>
                    <h1 className="text-4xl font-black text-slate-800 uppercase italic tracking-tighter">Календарь</h1>
                    <p className="text-slate-400 font-bold text-xs uppercase tracking-widest mt-1">
                        {monthNames[currentDate.getMonth()]} {currentDate.getFullYear()}
                    </p>
                </div>
                <div className="flex gap-2">
                    <button
                        onClick={() => setCurrentDate(new Date(currentDate.setMonth(currentDate.getMonth() - 1)))}
                        className="p-3 bg-white border border-slate-100 rounded-2xl hover:bg-slate-50"
                    >
                        ←
                    </button>
                    <button
                        onClick={() => setCurrentDate(new Date(currentDate.setMonth(currentDate.getMonth() + 1)))}
                        className="p-3 bg-white border border-slate-100 rounded-2xl hover:bg-slate-50"
                    >
                        →
                    </button>
                </div>
            </header>

            <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
                {/* Сетка календаря */}
                <div className="lg:col-span-2 bg-white rounded-[32px] border border-slate-100 p-6 shadow-sm">
                    <div className="grid grid-cols-7 mb-4">
                        {['Пн', 'Вт', 'Ср', 'Чт', 'Пт', 'Сб', 'Вс'].map(d => (
                            <div key={d} className="text-center text-[10px] font-black text-slate-300 uppercase">{d}</div>
                        ))}
                    </div>
                    <div className="grid grid-cols-7 gap-2">
                        {days.map((day, idx) => (
                            <div
                                key={idx}
                                className={`
                                    h-16 md:h-24 rounded-2xl flex flex-col items-center justify-center relative transition-all
                                    ${day ? 'bg-slate-50 hover:bg-[#2D9396]/10 cursor-pointer' : 'bg-transparent'}
                                    ${day === new Date().getDate() && currentDate.getMonth() === new Date().getMonth() ? 'ring-2 ring-[#2D9396] bg-white' : ''}
                                `}
                            >
                                <span className={`font-bold ${day === new Date().getDate() ? 'text-[#2D9396]' : 'text-slate-700'}`}>
                                    {day}
                                </span>
                                {/* Маркер события */}
                                {day % 7 === 0 && day !== null && (
                                    <div className="absolute bottom-2 w-1.5 h-1.5 bg-[#FF7A50] rounded-full"></div>
                                )}
                            </div>
                        ))}
                    </div>
                </div>

                {/* Список ближайших событий */}
                <div className="space-y-6">
                    <h2 className="text-sm font-black text-slate-400 uppercase tracking-widest">Ближайшие события</h2>
                    <div className="space-y-4">
                        {events.length > 0 ? events.map((event, i) => (
                            <div key={i} className="bg-white p-6 rounded-[28px] border border-slate-100 shadow-sm hover:translate-x-1 transition-transform cursor-pointer">
                                <div className="flex justify-between items-start mb-2">
                                    <span className="text-[10px] font-black text-[#2D9396] uppercase bg-[#2D9396]/10 px-3 py-1 rounded-lg">
                                        {event.type || 'Событие'}
                                    </span>
                                    <span className="text-slate-400 text-[10px] font-bold">{event.time}</span>
                                </div>
                                <h4 className="font-bold text-slate-800">{event.title}</h4>
                                <p className="text-slate-400 text-xs mt-1">{event.location || 'Онлайн'}</p>
                            </div>
                        )) : (
                            <div className="p-8 bg-slate-50 rounded-[32px] border-2 border-dashed border-slate-100 text-center">
                                <p className="text-slate-400 font-bold text-[10px] uppercase">Событий не найдено</p>
                            </div>
                        )}
                    </div>
                </div>
            </div>
        </div>
    );
};

export default CalendarPage;