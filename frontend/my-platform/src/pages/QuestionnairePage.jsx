import React, { useState, useEffect } from 'react';
import { api } from '../api';

const QuestionnairePage = () => {
    const [template, setTemplate] = useState(null);
    const [answers, setAnswers] = useState({});
    const [loading, setLoading] = useState(true);
    const [submitting, setSubmitting] = useState(false);

    useEffect(() => {
        const loadData = async () => {
            try {
                // Загружаем шаблон вопросов
                const res = await api.questionnaires.getTemplate();
                setTemplate(res.data);

                // Инициализируем пустые ответы на основе полученного шаблона
                const initialAnswers = {};
                res.data.sections.forEach(section => {
                    section.fields.forEach(field => {
                        initialAnswers[field.id] = "";
                    });
                });
                setAnswers(initialAnswers);
            } catch (err) {
                console.error("Ошибка при получении шаблона анкеты:", err);
            } finally {
                setLoading(false);
            }
        };
        loadData();
    }, []);

    const handleInputChange = (fieldId, value) => {
        setAnswers(prev => ({
            ...prev,
            [fieldId]: value
        }));
    };

    const handleSubmit = async (e) => {
        e.preventDefault();
        setSubmitting(true);
        try {
            // Формируем массив объектов согласно документации API
            const formattedAnswers = Object.entries(answers).map(([field_id, value]) => ({
                field_id,
                value
            }));

            await api.questionnaires.submit(formattedAnswers);
            alert("Анкета успешно отправлена!");
        } catch (err) {
            console.error("Ошибка отправки:", err);
            alert("Не удалось отправить анкету. Проверьте заполнение полей.");
        } finally {
            setSubmitting(false);
        }
    };

    if (loading) return (
        <div className="p-10 flex items-center justify-center">
            <div className="text-[#2D9396] font-black animate-pulse uppercase tracking-widest">
                Загрузка вопросов...
            </div>
        </div>
    );

    return (
        <div className="max-w-4xl mx-auto p-4 md:p-8 animate-in fade-in slide-in-from-bottom-4 duration-700">
            <header className="mb-10 text-center md:text-left">
                <h1 className="text-4xl font-black text-slate-800 uppercase italic tracking-tighter">
                    Анкета студента
                </h1>
                <p className="text-slate-400 mt-2 font-medium">Пожалуйста, заполните все обязательные поля для обновления профиля.</p>
            </header>

            <form onSubmit={handleSubmit} className="space-y-8">
                {template?.sections?.map((section) => (
                    <div key={section.id} className="bg-white rounded-[32px] border border-slate-100 shadow-sm p-8">
                        <h2 className="text-xs font-black text-[#2D9396] uppercase tracking-[0.2em] mb-8 border-b border-slate-50 pb-4">
                            {section.title}
                        </h2>

                        <div className="grid grid-cols-1 gap-8">
                            {section.fields.map((field) => (
                                <div key={field.id} className="flex flex-col gap-2">
                                    <label className="text-sm font-bold text-slate-700 ml-1">
                                        {field.label} {field.required && <span className="text-red-400">*</span>}
                                    </label>

                                    {field.type === 'textarea' ? (
                                        <textarea
                                            required={field.required}
                                            value={answers[field.id]}
                                            onChange={(e) => handleInputChange(field.id, e.target.value)}
                                            className="w-full bg-slate-50 border-none rounded-2xl p-4 min-h-[120px] focus:ring-2 focus:ring-[#2D9396] transition-all"
                                            placeholder="Введите ваш ответ..."
                                        />
                                    ) : (
                                        <input
                                            type={field.type || 'text'}
                                            required={field.required}
                                            value={answers[field.id]}
                                            onChange={(e) => handleInputChange(field.id, e.target.value)}
                                            className="w-full bg-slate-50 border-none rounded-2xl p-4 h-14 focus:ring-2 focus:ring-[#2D9396] transition-all"
                                            placeholder={field.placeholder || "..."}
                                        />
                                    )}
                                </div>
                            ))}
                        </div>
                    </div>
                ))}

                <div className="flex justify-center md:justify-end pt-4">
                    <button
                        type="submit"
                        disabled={submitting}
                        className={`
                            px-12 py-5 rounded-[24px] font-black uppercase text-xs tracking-widest transition-all shadow-lg
                            ${submitting
                            ? 'bg-slate-200 text-slate-400 cursor-not-allowed'
                            : 'bg-[#2D9396] text-white hover:bg-[#257d7f] hover:scale-[1.02] active:scale-95 shadow-[#2D9396]/20'}
                        `}
                    >
                        {submitting ? 'Отправка...' : 'Сохранить анкету'}
                    </button>
                </div>
            </form>
        </div>
    );
};

export default QuestionnairePage;