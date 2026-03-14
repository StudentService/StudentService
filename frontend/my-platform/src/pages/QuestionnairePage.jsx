import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { api } from '../api';

const QuestionnairePage = () => {
    const [template, setTemplate] = useState(null);
    const [answers, setAnswers] = useState({}); // Храним ответы тут: { field_id: "значение" }
    const [status, setStatus] = useState('loading'); // loading, ready, submitting, success
    const navigate = useNavigate();

    // 1. Загружаем вопросы при входе
    useEffect(() => {
        api.questionnaires.getTemplate()
            .then(res => {
                setTemplate(res.data);
                setStatus('ready');
            })
            .catch(err => {
                console.error("Ошибка загрузки анкеты", err);
                setStatus('error');
            });
    }, []);

    // 2. Обработка ввода
    const handleChange = (id, value) => {
        setAnswers(prev => ({ ...prev, [id]: value }));
    };

    // 3. Отправка формы
    const handleSubmit = async (e) => {
        e.preventDefault();
        setStatus('submitting');

        // Преобразуем объект {id: val} в массив [{field_id: id, value: val}] для API
        const payload = Object.entries(answers).map(([id, val]) => ({
            field_id: id,
            value: val
        }));

        try {
            await api.questionnaires.submit(payload);
            setStatus('success');
            setTimeout(() => navigate('/dashboard'), 2000); // Уводим на дашборд через 2 сек
        } catch (err) {
            alert("Ошибка при сохранении: " + (err.response?.data?.message || "Сервер недоступен"));
            setStatus('ready');
        }
    };

    if (status === 'loading') return <div className="p-10 text-[#2D9396] font-black animate-pulse">ЧИТАЕМ ВОПРОСЫ...</div>;

    if (status === 'success') return (
        <div className="p-10 text-center animate-bounce">
            <h2 className="text-2xl font-bold text-green-500 uppercase">Данные сохранены!</h2>
            <p className="text-slate-400">Перенаправляем в личный кабинет...</p>
        </div>
    );

    return (
        <div className="max-w-3xl mx-auto p-6">
            <form onSubmit={handleSubmit} className="space-y-6">
                {template?.sections?.map(section => (
                    <div key={section.id} className="bg-white p-8 rounded-[40px] shadow-sm border border-slate-100">
                        <h2 className="text-[#2D9396] font-black uppercase text-xs tracking-widest mb-6">{section.title}</h2>
                        <div className="space-y-4">
                            {section.fields.map(field => (
                                <div key={field.id} className="flex flex-col gap-2">
                                    <label className="text-xs font-bold text-slate-500 ml-2">{field.label}</label>
                                    <input
                                        type={field.type}
                                        required={field.required}
                                        className="bg-slate-50 border-none rounded-2xl p-4 focus:ring-2 focus:ring-[#2D9396]"
                                        onChange={(e) => handleChange(field.id, e.target.value)}
                                    />
                                </div>
                            ))}
                        </div>
                    </div>
                ))}

                <button
                    type="submit"
                    disabled={status === 'submitting'}
                    className="w-full bg-[#2D9396] text-white py-5 rounded-[24px] font-black uppercase shadow-lg shadow-[#2D9396]/20"
                >
                    {status === 'submitting' ? 'СОХРАНЯЕМ...' : 'ОТПРАВИТЬ АНКЕТУ'}
                </button>
            </form>
        </div>
    );
};

export default QuestionnairePage;