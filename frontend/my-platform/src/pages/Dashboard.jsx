import React, { useState, useEffect } from 'react';
import { useAuth } from '../contexts/AuthContext';
import { api } from '../api';

const Dashboard = () => {
    const { user } = useAuth();
    const [data, setData] = useState(null);
    const [loading, setLoading] = useState(true);

    const isTeacher = user?.role === 'teacher';

    useEffect(() => {
        const loadDashboard = async () => {
            try {
                let res;
                if (isTeacher) {
                    res = await api.dashboard.getTeacher();
                } else {
                    res = await api.dashboard.getStudent();
                }
                setData(res.data);
            } catch (err) {
                console.error("Dashboard load error", err);
            } finally {
                setLoading(false);
            }
        };
        loadDashboard();
    }, [isTeacher]);

    if (loading) return <div className="p-10 font-black text-[#2D9396] animate-pulse">ЗАГРУЗКА...</div>;

    // Дашборд преподавателя
    if (isTeacher) {
        const { groups, students_count, activities } = data || {};
        return (
            <div className="p-4 md:p-8 space-y-8 animate-in fade-in duration-500">
                {/* Баннер */}
                <div className="bg-[#2D9396] rounded-[32px] p-10 text-white shadow-lg relative overflow-hidden">
                    <h1 className="text-4xl font-black mb-2 italic uppercase tracking-tighter">
                        Привет, {user?.name?.split(' ')[0] || 'Преподаватель'}!
                    </h1>
                    <p className="text-lg opacity-90">У вас {groups?.length || 0} групп и {students_count || 0} студентов.</p>
                </div>

                {/* Метрики */}
                <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
                    <MetricCard title="Группы" value={groups?.length || 0} icon="👥" subtitle="Ваше руководство" />
                    <MetricCard title="Студенты" value={students_count || 0} icon="🎓" subtitle="Всего" />
                    <MetricCard title="Активности" value={activities?.length || 0} icon="⚡" subtitle="Проведено" />
                </div>

                {/* Список групп */}
                <div className="bg-white rounded-[32px] border border-slate-100 p-8 shadow-sm">
                    <h2 className="text-xl font-black text-slate-800 uppercase italic mb-8">Ваши группы</h2>
                    <div className="space-y-4">
                        {groups?.map(group => (
                            <div key={group.id} className="flex justify-between items-center p-4 bg-slate-50 rounded-2xl">
                                <div>
                                    <h4 className="font-bold text-slate-800">{group.name}</h4>
                                    <p className="text-slate-400 text-[10px] font-bold uppercase">{group.faculty}</p>
                                </div>
                                <span className="text-[#2D9396] font-black">{group.students_count} студ.</span>
                            </div>
                        ))}
                    </div>
                </div>
            </div>
        );
    }

    // Дашборд студента
    const { student_info, statistics, active_challenges } = data || {};

    return (
        <div className="p-4 md:p-8 space-y-8 animate-in fade-in duration-500">
            {/* Баннер */}
            <div className="bg-[#2D9396] rounded-[32px] p-10 text-white shadow-lg relative overflow-hidden">
                <h1 className="text-4xl font-black mb-2 italic uppercase tracking-tighter">
                    Привет, {student_info?.name?.split(' ')[0] || 'Студент'}!
                </h1>
                <p className="text-lg opacity-90">У тебя {active_challenges?.length || 0} активных вызовов.</p>
            </div>

            {/* Метрики */}
            <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
                <MetricCard title="Ядро" subtitle={student_info?.group_name} icon="🏆" badge="Текущий" />
                <MetricCard title="Прогресс" value={statistics?.total_points} icon="🎯" progress={68} />
                <MetricCard title={statistics?.average_grade || 0} subtitle="Средний балл" icon="📈" />
            </div>

            {/* Активные вызовы */}
            <div className="bg-white rounded-[32px] border border-slate-100 p-8 shadow-sm">
                <h2 className="text-xl font-black text-slate-800 uppercase italic mb-8">Активные вызовы</h2>
                <div className="space-y-6">
                    {active_challenges?.map(ch => (
                        <div key={ch.id} className="group">
                            <div className="flex justify-between items-end mb-2">
                                <div>
                                    <h4 className="font-bold text-slate-800">{ch.title}</h4>
                                    <span className="text-slate-400 text-[10px] font-bold uppercase">Осталось дней: {ch.days_left}</span>
                                </div>
                                <span className="font-black text-[#2D9396]">{ch.progress}%</span>
                            </div>
                            <div className="h-2 bg-slate-50 rounded-full overflow-hidden">
                                <div className="h-full bg-[#2D9396]" style={{ width: `${ch.progress}%` }}></div>
                            </div>
                        </div>
                    ))}
                </div>
            </div>
        </div>
    );
};

const MetricCard = ({ title, subtitle, icon, progress, badge, value }) => (
    <div className="bg-white p-8 rounded-[32px] border border-slate-100 shadow-sm relative">
        <div className="flex justify-between items-start mb-4">
            <div className="text-2xl">{icon}</div>
            {badge && <span className="bg-[#2D9396] text-white text-[10px] px-3 py-1 rounded-lg font-black uppercase">{badge}</span>}
            {value !== undefined && <span className="text-3xl font-black text-[#2D9396]">{value}</span>}
        </div>
        <h3 className="text-3xl font-black text-slate-800">{title}</h3>
        {subtitle && <p className="text-slate-400 text-[10px] font-black uppercase mt-1">{subtitle}</p>}
        {progress && (
            <div className="mt-4 h-1.5 bg-slate-50 rounded-full overflow-hidden">
                <div className="h-full bg-[#FF7A50]" style={{ width: `${progress}%` }}></div>
            </div>
        )}
    </div>
);

export default Dashboard;