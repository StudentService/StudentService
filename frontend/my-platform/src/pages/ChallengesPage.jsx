import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { api } from '../api';

const ChallengePage = () => {
    const navigate = useNavigate();
    const [loading, setLoading] = useState(false);

    const [formData, setFormData] = useState({
        title: '',
        description: '', // Для req.Description
        goal: '',        // Для req.Goal
        startDate: new Date().toISOString().split('T')[0], // Сегодня по умолчанию
        endDate: ''      // Для req.EndDate
    });

    const handleSubmit = async (e) => {
        e.preventDefault();
        setLoading(true);

        try {
            // 1. Берем даты из формы
            const start = new Date(formData.startDate);
            const end = new Date(formData.endDate);

            // 2. Валидация на фронте (чтобы не гонять лишний запрос)
            if (end <= start) {
                alert("Дата окончания должна быть позже даты начала!");
                setLoading(false);
                return;
            }

            // 3. Формируем Payload СТРОГО по твоим Go-тегам json:"..."
            const payload = {
                title: formData.title.trim(),
                description: formData.description.trim(), // может быть пустым
                goal: formData.goal.trim(),
                start_date: start.toISOString(), // Снейк-кейс!
                end_date: end.toISOString()      // Снейк-кейс!
            };

            console.log("Отправляем идеальный Payload:", payload);

            await api.challenges.create(payload);

            alert("Успех! Челлендж создан.");
            navigate('/challenges');

        } catch (err) {
            console.error("Ошибка сервера:", err.response?.data);
            // Выводим конкретную ошибку валидации из Go
            const errorMsg = err.response?.data?.error || "Ошибка сервера";
            alert(`Сервер не принял данные: ${errorMsg}`);
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className="max-w-2xl mx-auto p-4 md:p-8 animate-in fade-in slide-in-from-bottom-4 duration-500">
            <header className="mb-8">
                <h1 className="text-4xl font-black text-slate-800 uppercase italic tracking-tighter">Новый вызов</h1>
                <p className="text-slate-400 font-bold text-[10px] uppercase tracking-widest mt-1">StudentService Challenge System</p>
            </header>

            <form onSubmit={handleSubmit} className="bg-white rounded-[40px] border border-slate-100 shadow-sm p-8 space-y-6">
                <div className="space-y-2">
                    <label className="text-xs font-black text-slate-400 uppercase ml-2">Заголовок</label>
                    <input
                        required
                        type="text"
                        className="w-full bg-slate-50 border-none rounded-2xl p-4 focus:ring-2 focus:ring-[#2D9396] transition-all"
                        placeholder="Название вызова"
                        onChange={(e) => setFormData({...formData, title: e.target.value})}
                    />
                </div>

                <div className="space-y-2">
                    <label className="text-xs font-black text-slate-400 uppercase ml-2">Краткое описание</label>
                    <input
                        required
                        type="text"
                        className="w-full bg-slate-50 border-none rounded-2xl p-4 focus:ring-2 focus:ring-[#2D9396]"
                        placeholder="О чем этот вызов?"
                        onChange={(e) => setFormData({...formData, description: e.target.value})}
                    />
                </div>

                <div className="space-y-2">
                    <label className="text-xs font-black text-slate-400 uppercase ml-2">Финальная цель (Goal)</label>
                    <textarea
                        required
                        className="w-full bg-slate-50 border-none rounded-2xl p-4 min-h-[100px] focus:ring-2 focus:ring-[#2D9396] transition-all"
                        placeholder="Чего конкретно ты хочешь достичь?"
                        onChange={(e) => setFormData({...formData, goal: e.target.value})}
                    />
                </div>

                <div className="grid grid-cols-2 gap-4">
                    <div className="space-y-2">
                        <label className="text-xs font-black text-slate-400 uppercase ml-2">Дата начала</label>
                        <input
                            required
                            type="date"
                            value={formData.startDate}
                            className="w-full bg-slate-50 border-none rounded-2xl p-4 focus:ring-2 focus:ring-[#2D9396]"
                            onChange={(e) => setFormData({...formData, startDate: e.target.value})}
                        />
                    </div>
                    <div className="space-y-2">
                        <label className="text-xs font-black text-slate-400 uppercase ml-2">Дата окончания</label>
                        <input
                            required
                            type="date"
                            className="w-full bg-slate-50 border-none rounded-2xl p-4 focus:ring-2 focus:ring-[#2D9396]"
                            onChange={(e) => setFormData({...formData, endDate: e.target.value})}
                        />
                    </div>
                </div>

                <button
                    type="submit"
                    disabled={loading}
                    className="w-full bg-[#2D9396] text-white py-5 rounded-[24px] font-black uppercase text-xs tracking-widest shadow-lg shadow-[#2D9396]/20 hover:scale-[1.01] active:scale-95 transition-all"
                >
                    {loading ? 'Создание...' : 'Опубликовать вызов'}
                </button>
            </form>
        </div>
    );
};

export default ChallengePage;