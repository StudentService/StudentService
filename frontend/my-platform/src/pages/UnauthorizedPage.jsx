import React from 'react';
import { Link } from 'react-router-dom';

const UnauthorizedPage = () => {
    return (
        <div className="min-h-screen bg-slate-50 flex items-center justify-center p-4">
            <div className="text-center">
                <h1 className="text-6xl font-black text-red-500 mb-4">403</h1>
                <h2 className="text-2xl font-bold text-slate-800 mb-4">Доступ запрещён</h2>
                <p className="text-slate-500 mb-8 max-w-md">
                    У вас недостаточно прав для просмотра этой страницы.
                    Обратитесь к администратору, если считаете, что это ошибка.
                </p>
                <Link
                    to="/dashboard"
                    className="inline-block bg-blue-600 hover:bg-blue-700 text-white font-bold py-3 px-8 rounded-2xl transition-all"
                >
                    На главную
                </Link>
            </div>
        </div>
    );
};

export default UnauthorizedPage;