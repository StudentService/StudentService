import React from 'react';

const ChallengeCard = ({ title, semester, status, current, total, progress }) => (
    <div className="bg-white rounded-3xl border border-slate-200 p-6 mb-4 hover:shadow-md transition-shadow group cursor-pointer">
        <div className="flex justify-between items-start mb-4">
            <div className="flex gap-4">
                <div className="w-12 h-12 bg-slate-50 rounded-full flex items-center justify-center border border-brand/20 text-brand">
                    🎯
                </div>
                <div>
                    <h3 className="text-xl font-bold text-slate-900 group-hover:text-brand transition-colors">{title}</h3>
                    <div className="flex gap-3 mt-1">
                        <span className="text-sm text-slate-500">Семестр: {semester} семестр</span>
                        <span className={`text-[10px] px-2 py-0.5 rounded-full font-bold uppercase ${
                            status === 'Начат' ? 'bg-emerald-100 text-emerald-600' : 'bg-blue-100 text-blue-600'
                        }`}>
              {status}
            </span>
                        <span className="text-sm text-slate-400 font-medium italic">Чекпоинты: {current}/{total}</span>
                    </div>
                </div>
            </div>
            <span className="text-slate-300">❯</span>
        </div>

        <div>
            <div className="flex justify-between text-xs font-bold text-slate-400 mb-2 uppercase tracking-tight">
                <span>Прогресс</span>
                <span className="text-brand">{progress}%</span>
            </div>
            <div className="w-full bg-slate-100 h-3 rounded-full overflow-hidden">
                <div
                    className="bg-brand h-full rounded-full transition-all duration-1000"
                    style={{ width: `${progress}%` }}
                ></div>
            </div>
        </div>
    </div>
);

const Challenges = () => {
    return (
        <div className="max-w-4xl">
            <div className="flex justify-between items-end mb-8">
                <div>
                    <h1 className="text-3xl font-black text-slate-900">Мои вызовы</h1>
                    <p className="text-slate-500 mt-1">Управляйте своими личными вызовами и отслеживайте прогресс</p>
                </div>
                <button className="bg-brand hover:bg-brand-dark text-white px-6 py-2.5 rounded-xl font-bold flex items-center gap-2 transition-all shadow-lg shadow-brand/20">
                    <span className="text-xl">+</span> Создать новый вызов
                </button>
            </div>

            <div className="bg-white rounded-2xl border border-slate-200 p-4 mb-6 flex items-center gap-4">
                <span className="text-slate-400">🔍</span>
                <label className="text-sm font-medium text-slate-600">Фильтр по семестру:</label>
                <select className="bg-slate-50 border border-slate-200 rounded-lg px-3 py-1 text-sm outline-none font-bold text-slate-700">
                    <option>Все семестры</option>
                    <option>6 семестр</option>
                </select>
            </div>

            <ChallengeCard
                title="Разработка веб-приложения"
                semester="6"
                status="В процессе"
                current="3"
                total="4"
                progress={75}
            />
            <ChallengeCard
                title="Исследовательский проект"
                semester="6"
                status="Начат"
                current="2"
                total="5"
                progress={45}
            />
        </div>
    );
};

export default Challenges;