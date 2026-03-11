import React from 'react';
import { Outlet, Link, useLocation, useNavigate } from 'react-router-dom';

const MainLayout = () => {
    const location = useLocation();
    const navigate = useNavigate();

    // Список пунктов меню согласно твоему макету
    const menuItems = [
        { name: 'Дашборд', path: '/dashboard', icon: '🏠' },
        { name: 'Профиль', path: '/profile', icon: '👤' },
        { name: 'Анкета', path: '/form', icon: '📄' },
        { name: 'Вызовы', path: '/challenges', icon: '🎯' },
        { name: 'Календарь', path: '/calendar', icon: '📅' },
        { name: 'Успеваемость', path: '/stats', icon: '📈' },
        { name: 'Активности', path: '/activities', icon: '⚡' },
        { name: 'Метрики', path: '/metrics', icon: '📊' },
    ];

    // Функция выхода
    const handleLogout = () => {
        localStorage.removeItem('access_token');
        // Используем replace: true, чтобы пользователь не мог вернуться назад кнопкой браузера
        navigate('/login', { replace: true });
    };

    return (
        <div className="min-h-screen bg-[#f8f9fa] flex flex-col font-sans">
            {/* --- ВЕРХНЯЯ ПАНЕЛЬ (HEADER) --- */}
            <header className="bg-white border-b border-slate-200 sticky top-0 z-20">
                <div className="max-w-[1440px] mx-auto flex items-center justify-between px-6 py-4">
                    <div className="flex items-center gap-10">
                        {/* Логотип */}
                        <div className="flex items-center gap-2 text-brand font-black text-xl uppercase tracking-widest">
                            <span className="bg-brand text-white px-1.5 py-0.5 rounded shadow-sm">▲</span>
                            Капитаны
                        </div>
                        {/* Горизонтальная навигация (как на скрине) */}
                        <nav className="hidden lg:flex gap-8 text-[13px] font-bold uppercase tracking-tight text-slate-400">
                            {menuItems.slice(0, 4).map((item) => (
                                <Link
                                    key={item.path}
                                    to={item.path}
                                    className={`transition-colors hover:text-brand ${location.pathname === item.path ? 'text-brand border-b-2 border-brand pb-1' : ''}`}
                                >
                                    {item.name}
                                </Link>
                            ))}
                        </nav>
                    </div>

                    {/* Правая часть шапки */}
                    <div className="flex items-center gap-3">
                        <div className="text-right hidden sm:block">
                            <p className="text-xs font-bold text-slate-800 italic uppercase">Администратор</p>
                        </div>
                        <div className="w-9 h-9 bg-slate-100 rounded-full border border-slate-200 flex items-center justify-center text-slate-400">
                            👤
                        </div>
                    </div>
                </div>
            </header>

            {/* --- ОСНОВНОЙ КОНТЕНТ --- */}
            <div className="flex max-w-[1440px] w-full mx-auto p-6 gap-8">

                {/* БОКОВОЕ МЕНЮ (SIDEBAR) */}
                <aside className="w-64 flex-shrink-0">
                    <div className="bg-white rounded-3xl border border-slate-200 p-4 sticky top-24 shadow-sm flex flex-col h-[calc(100vh-120px)]">

                        {/* Пункты навигации */}
                        <nav className="space-y-1 flex-1 overflow-y-auto pr-1">
                            {menuItems.map((item) => {
                                const isActive = location.pathname === item.path;
                                return (
                                    <Link
                                        key={item.path}
                                        to={item.path}
                                        className={`flex items-center gap-4 px-4 py-3.5 rounded-2xl text-sm font-bold transition-all ${
                                            isActive
                                                ? 'bg-brand text-white shadow-lg shadow-brand/30 scale-[1.02]'
                                                : 'text-slate-400 hover:bg-slate-50 hover:text-slate-600'
                                        }`}
                                    >
                                        <span className="text-lg">{item.icon}</span>
                                        {item.name}
                                    </Link>
                                );
                            })}
                        </nav>

                        {/* КНОПКА ВЫХОДА */}
                        <div className="mt-auto pt-4 border-t border-slate-100">
                            <button
                                onClick={handleLogout}
                                className="w-full flex items-center gap-4 px-4 py-3.5 rounded-2xl text-sm font-bold text-red-400 hover:bg-red-50 hover:text-red-500 transition-all cursor-pointer"
                            >
                                <span className="text-lg">🚪</span>
                                Выйти
                            </button>
                        </div>
                    </div>
                </aside>

                {/* ЗОНА ВЫВОДА СТРАНИЦ */}
                <main className="flex-1 min-w-0">
                    <Outlet />
                </main>

            </div>
        </div>
    );
};

export default MainLayout;