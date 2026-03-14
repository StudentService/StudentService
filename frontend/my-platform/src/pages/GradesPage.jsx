import React, { useState, useEffect } from 'react';
import { api } from '../api';

const GradesPage = () => {
    const [summary, setSummary] = useState(null);
    const [history, setHistory] = useState([]);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        const fetchGrades = async () => {
            try {
                const [summaryRes, historyRes] = await Promise.all([
                    api.grades.getSummary(),
                    api.grades.getAll()
                ]);
                setSummary(summaryRes.data);
                setHistory(historyRes.data);
            } catch (err) {
                console.error("Ошибка загрузки успеваемости:", err);
            } finally {
                setLoading(false);
            }
        };
        fetchGrades();
    }, []);

    if (loading) return <div className="p-10 font-black text-[#2D9396] animate-pulse uppercase">Анализ успеваемости...</div>;

    return (
        <div className="p-4 md:p-8 space-y-8 animate-in fade-in duration-700">
            <header>
                <h1 className="text-4xl font-black text-slate-800 uppercase italic tracking-tighter">Успеваемость</h1>
                <div className="flex items-center gap-4 mt-2">
                    <span className="text-[10px] font-black bg-[#FF7A50] text-white px-3 py-1 rounded-lg uppercase">
                        GPA: {summary?.overall_average?.toFixed(2) || '0.00'}
                    </span>
                    <span className="text-slate-400 text-[10px] font-black uppercase tracking-widest">
                        Всего кредитов: {summary?.total_credits || 0}
                    </span>
                </div>
            </header>

            {/* Сводка по курсам */}
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                {summary?.courses?.map((course) => (
                    <div key={course.course_id} className="bg-white p-6 rounded-[32px] border border-slate-100 shadow-sm hover:shadow-md transition-shadow">
                        <div className="flex justify-between items-start mb-4">
                            <div className="w-10 h-10 bg-slate-50 rounded-xl flex items-center justify-center text-lg">📚</div>
                            <div className="text-right">
                                <p className="text-[10px] font-black text-slate-300 uppercase">Средний балл</p>
                                <p className="text-xl font-black text-[#2D9396]">{course.average?.toFixed(1)}</p>
                            </div>
                        </div>
                        <h3 className="font-bold text-slate-800 mb-1">{course.course_name}</h3>
                        <p className="text-slate-400 text-[10px] font-bold uppercase tracking-tighter">
                            Оценок: {course.grades_count} • Кредиты: {course.total_credits}
                        </p>
                        <div className="mt-4 h-1.5 bg-slate-50 rounded-full overflow-hidden">
                            <div
                                className="h-full bg-[#2D9396]"
                                style={{ width: `${(course.average / 5) * 100}%` }}
                            ></div>
                        </div>
                    </div>
                ))}
            </div>

            {/* История оценок */}
            <div className="bg-white rounded-[32px] border border-slate-100 shadow-sm overflow-hidden">
                <div className="p-8 border-b border-slate-50">
                    <h2 className="text-sm font-black text-slate-400 uppercase tracking-widest">История последних оценок</h2>
                </div>
                <div className="overflow-x-auto">
                    <table className="w-full text-left border-collapse">
                        <thead>
                        <tr className="bg-slate-50/50">
                            <th className="p-6 text-[10px] font-black text-slate-400 uppercase">Дата</th>
                            <th className="p-6 text-[10px] font-black text-slate-400 uppercase">Курс</th>
                            <th className="p-6 text-[10px] font-black text-slate-400 uppercase">Тип</th>
                            <th className="p-6 text-[10px] font-black text-slate-400 uppercase text-center">Балл</th>
                            <th className="p-6 text-[10px] font-black text-slate-400 uppercase">Комментарий</th>
                        </tr>
                        </thead>
                        <tbody className="divide-y divide-slate-50">
                        {history.map((grade) => (
                            <tr key={grade.id} className="hover:bg-slate-50/30 transition-colors">
                                <td className="p-6 text-sm text-slate-500 font-medium">
                                    {new Date(grade.date).toLocaleDateString()}
                                </td>
                                <td className="p-6 font-bold text-slate-800">{grade.course_name}</td>
                                <td className="p-6">
                                        <span className="text-[9px] font-black uppercase px-2 py-1 bg-slate-100 text-slate-500 rounded">
                                            {grade.type}
                                        </span>
                                </td>
                                <td className="p-6 text-center">
                                        <span className={`
                                            text-lg font-black 
                                            ${grade.value >= 4 ? 'text-[#2D9396]' : grade.value === 3 ? 'text-orange-400' : 'text-red-400'}
                                        `}>
                                            {grade.value}
                                        </span>
                                    <span className="text-slate-300 text-[10px] ml-1">/ {grade.max_value}</span>
                                </td>
                                <td className="p-6 text-sm text-slate-400 italic">
                                    {grade.comment || '—'}
                                </td>
                            </tr>
                        ))}
                        </tbody>
                    </table>
                </div>
                {history.length === 0 && (
                    <div className="p-20 text-center text-slate-300 font-bold uppercase text-xs">Оценок пока нет</div>
                )}
            </div>
        </div>
    );
};

export default GradesPage;