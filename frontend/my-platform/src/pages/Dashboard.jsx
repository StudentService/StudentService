import React from 'react';

const Dashboard = () => {
    // Данные для карточек (в будущем придут из API)
    const stats = [
        { id: 1, label: 'Всего студентов', value: '1,284', change: '+12%', color: 'text-blue-600', bg: 'bg-blue-50' },
        { id: 2, label: 'Активные (Ядро)', value: '312', change: '+5%', color: 'text-emerald-600', bg: 'bg-emerald-50' },
        { id: 3, label: 'Мероприятия', value: '18', change: 'в этом мес.', color: 'text-purple-600', bg: 'bg-purple-50' },
        { id: 4, label: 'Средний балл', value: '4.2', change: '+0.3', color: 'text-amber-600', bg: 'bg-amber-50' },
    ];

    return (
        <div className="space-y-8">
            {/* Заголовок */}
            <div>
                <h1 className="text-2xl font-bold text-slate-800">Общая аналитика</h1>
                <p className="text-slate-500 text-sm">Добро пожаловать в систему управления Центром Развития</p>
            </div>

            {/* Сетка карточек */}
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
                {stats.map((stat) => (
                    <div key={stat.id} className="bg-white p-6 rounded-3xl border border-slate-100 shadow-sm hover:shadow-md transition-shadow">
                        <div className={`w-12 h-12 ${stat.bg} ${stat.color} rounded-2xl flex items-center justify-center mb-4 text-xl font-bold`}>
                            {stat.label[0]}
                        </div>
                        <p className="text-slate-500 text-sm font-medium">{stat.label}</p>
                        <div className="flex items-end justify-between mt-1">
                            <p className="text-3xl font-bold text-slate-800">{stat.value}</p>
                            <span className="text-emerald-500 text-xs font-bold mb-1">{stat.change}</span>
                        </div>
                    </div>
                ))}
            </div>

            {/* Секция с графиками (Визуализация из ТЗ) */}
            <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
                {/* Большая карточка: Активность */}
                <div className="lg:col-span-2 bg-white p-6 rounded-3xl border border-slate-100 shadow-sm">
                    <div className="flex justify-between items-center mb-6">
                        <h3 className="font-bold text-slate-800">Динамика набора баллов</h3>
                        <select className="text-sm bg-slate-50 border-none rounded-lg p-1 outline-none">
                            <option>За неделю</option>
                            <option>За месяц</option>
                        </select>
                    </div>
                    <div className="h-64 bg-slate-50 rounded-2xl flex items-center justify-center border-2 border-dashed border-slate-200">
                        <span className="text-slate-400 italic">Здесь будет график LineChart (Chart.js / Recharts)</span>
                    </div>
                </div>

                {/* Малая карточка: Сегментация (Ядро/Актив/Аура) */}
                <div className="bg-white p-6 rounded-3xl border border-slate-100 shadow-sm text-center">
                    <h3 className="font-bold text-slate-800 mb-6 text-left">Сегментация</h3>
                    <div className="relative inline-flex items-center justify-center mb-6">
                        {/* Имитация кругового графика через CSS */}
                        <div className="w-40 h-40 rounded-full border-[12px] border-blue-600 border-t-emerald-500 border-l-amber-400 rotate-45"></div>
                        <div className="absolute flex flex-col items-center">
                            <span className="text-2xl font-bold text-slate-800">100%</span>
                            <span className="text-[10px] text-slate-400 uppercase font-bold">Студентов</span>
                        </div>
                    </div>
                    <div className="space-y-2 text-left">
                        <div className="flex justify-between text-sm italic">
                            <span className="flex items-center gap-2"><div className="w-2 h-2 rounded-full bg-emerald-500"></div> Ядро</span>
                            <span className="font-bold">24%</span>
                        </div>
                        <div className="flex justify-between text-sm italic">
                            <span className="flex items-center gap-2"><div className="w-2 h-2 rounded-full bg-blue-600"></div> Актив</span>
                            <span className="font-bold">56%</span>
                        </div>
                        <div className="flex justify-between text-sm italic">
                            <span className="flex items-center gap-2"><div className="w-2 h-2 rounded-full bg-amber-400"></div> Аура</span>
                            <span className="font-bold">20%</span>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    );
};

export default Dashboard;