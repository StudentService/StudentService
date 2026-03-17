import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { api } from '../api';

const QuestionnairePage = () => {
    const [template, setTemplate] = useState(null);
    const [answers, setAnswers] = useState({});
    const [status, setStatus] = useState('loading'); // loading, ready, submitting, success, error
    const [existingQuestionnaire, setExistingQuestionnaire] = useState(null);
    const navigate = useNavigate();

    // 1. Загружаем шаблон и существующую анкету
    useEffect(() => {
        const loadData = async () => {
            try {
                const [templateRes, myRes] = await Promise.all([
                    api.questionnaires.getTemplate(),
                    api.questionnaires.getMy().catch(() => null) // Если анкеты нет - игнорируем ошибку
                ]);

                setTemplate(templateRes.data);

                // Если есть сохранённая анкета, подгружаем ответы
                if (myRes?.data) {
                    setExistingQuestionnaire(myRes.data);
                    setAnswers(myRes.data.answers || {});
                    setStatus('edit'); // Режим редактирования
                } else {
                    setStatus('ready');
                }
            } catch (err) {
                console.error("Ошибка загрузки анкеты", err);
                setStatus('error');
            }
        };
        loadData();
    }, []);

    // 2. Обработка ввода
    const handleChange = (fieldId, value) => {
        setAnswers(prev => ({ ...prev, [fieldId]: value }));
    };

    // 3. Сохранение черновика
    const handleSaveDraft = async () => {
        try {
            await api.questionnaires.saveDraft(answers);
            alert('Черновик сохранён');
        } catch (err) {
            alert('Ошибка при сохранении черновика');
        }
    };

    // 4. Отправка анкеты
    const handleSubmit = async (e) => {
        e.preventDefault();
        setStatus('submitting');

        try {
            await api.questionnaires.submit(answers);
            setStatus('success');
            setTimeout(() => navigate('/dashboard'), 2000);
        } catch (err) {
            alert("Ошибка при сохранении: " + (err.response?.data?.error || "Сервер недоступен"));
            setStatus('ready');
        }
    };

    if (status === 'loading') return (
        <div className="p-10 text-center">
            <div className="inline-block animate-spin text-4xl mb-4">📋</div>
            <p className="text-[#2D9396] font-black animate-pulse">ЗАГРУЗКА АНКЕТЫ...</p>
        </div>
    );

    if (status === 'success') return (
        <div className="p-10 text-center animate-bounce">
            <div className="text-6xl mb-4">✅</div>
            <h2 className="text-2xl font-bold text-green-500 uppercase">Анкета отправлена!</h2>
            <p className="text-slate-400 mt-2">Перенаправляем в личный кабинет...</p>
        </div>
    );

    if (status === 'error') return (
        <div className="p-10 text-center">
            <div className="text-6xl mb-4">😢</div>
            <h2 className="text-2xl font-bold text-red-400 uppercase">Ошибка загрузки</h2>
            <button
                onClick={() => window.location.reload()}
                className="mt-6 px-8 py-4 bg-[#2D9396] text-white rounded-2xl font-black uppercase text-xs"
            >
                Попробовать снова
            </button>
        </div>
    );

    return (
        <div className="p-4 md:p-8 max-w-4xl mx-auto animate-in fade-in duration-500">
            {/* Заголовок */}
            <header className="mb-8">
                <h1 className="text-4xl font-black text-slate-800 uppercase italic tracking-tighter">
                    {existingQuestionnaire ? 'Редактирование анкеты' : 'Анкета студента'}
                </h1>
                <p className="text-slate-400 font-bold text-[10px] uppercase tracking-[0.2em] mt-1">
                    {existingQuestionnaire ? `Статус: ${existingQuestionnaire.status}` : 'Заполните все обязательные поля'}
                </p>
                {existingQuestionnaire?.status === 'approved' && (
                    <div className="mt-4 p-4 bg-green-50 text-green-600 rounded-2xl border border-green-100">
                        ✓ Анкета одобрена {existingQuestionnaire.reviewed_at ? new Date(existingQuestionnaire.reviewed_at).toLocaleDateString() : ''}
                    </div>
                )}
            </header>

            <form onSubmit={handleSubmit} className="space-y-8">
                {/* Динамические поля из шаблона */}
                {template?.fields?.map((field, idx) => (
                    <div key={field.id || idx} className="bg-white p-8 rounded-[40px] shadow-sm border border-slate-100">
                        <div className="flex items-center gap-2 mb-4">
                            <h2 className="text-sm font-black text-[#2D9396] uppercase tracking-widest">
                                {field.label || field.name}
                            </h2>
                            {field.required && (
                                <span className="text-[8px] font-black bg-red-50 text-red-400 px-2 py-1 rounded uppercase">
                                    обязательно
                                </span>
                            )}
                        </div>

                        {field.type === 'textarea' ? (
                            <textarea
                                required={field.required}
                                value={answers[field.id] || ''}
                                onChange={(e) => handleChange(field.id, e.target.value)}
                                className="w-full bg-slate-50 border-none rounded-2xl p-5 min-h-[120px] focus:ring-2 focus:ring-[#2D9396] transition-all"
                                placeholder={field.placeholder || ''}
                                disabled={existingQuestionnaire?.status === 'approved'}
                            />
                        ) : field.type === 'select' ? (
                            <select
                                required={field.required}
                                value={answers[field.id] || ''}
                                onChange={(e) => handleChange(field.id, e.target.value)}
                                className="w-full bg-slate-50 border-none rounded-2xl p-5 focus:ring-2 focus:ring-[#2D9396] transition-all"
                                disabled={existingQuestionnaire?.status === 'approved'}
                            >
                                <option value="">Выберите...</option>
                                {field.options?.map(opt => (
                                    <option key={opt.value} value={opt.value}>{opt.label}</option>
                                ))}
                            </select>
                        ) : (
                            <input
                                type={field.type || 'text'}
                                required={field.required}
                                value={answers[field.id] || ''}
                                onChange={(e) => handleChange(field.id, e.target.value)}
                                className="w-full bg-slate-50 border-none rounded-2xl p-5 focus:ring-2 focus:ring-[#2D9396] transition-all"
                                placeholder={field.placeholder || ''}
                                disabled={existingQuestionnaire?.status === 'approved'}
                            />
                        )}

                        {field.description && (
                            <p className="text-slate-400 text-xs mt-2 italic">{field.description}</p>
                        )}
                    </div>
                ))}

                {/* Кнопки действий */}
                {existingQuestionnaire?.status !== 'approved' && (
                    <div className="flex gap-4">
                        <button
                            type="submit"
                            disabled={status === 'submitting'}
                            className="flex-1 bg-[#2D9396] text-white py-6 rounded-[24px] font-black uppercase text-xs tracking-widest shadow-lg shadow-[#2D9396]/20 hover:scale-[1.01] active:scale-95 transition-all disabled:opacity-50"
                        >
                            {status === 'submitting' ? 'ОТПРАВКА...' : existingQuestionnaire ? 'ОБНОВИТЬ АНКЕТУ' : 'ОТПРАВИТЬ АНКЕТУ'}
                        </button>

                        <button
                            type="button"
                            onClick={handleSaveDraft}
                            className="px-8 py-6 bg-slate-50 hover:bg-slate-100 rounded-[24px] font-black uppercase text-xs tracking-widest transition-all"
                        >
                            СОХРАНИТЬ ЧЕРНОВИК
                        </button>
                    </div>
                )}

                {existingQuestionnaire?.status === 'approved' && (
                    <div className="text-center p-8 bg-slate-50 rounded-[40px]">
                        <p className="text-slate-500">Анкета уже одобрена и не может быть изменена</p>
                    </div>
                )}
            </form>
        </div>
    );
};

export default QuestionnairePage;