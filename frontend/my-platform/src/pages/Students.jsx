import React, { useState, useEffect } from 'react';
import { api } from '../api/index.js';

const Students = () => {
    const [students, setStudents] = useState([]);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        api.getStudents()
            .then(res => setStudents(res.data.data || res.data))
            .catch(err => console.error("Ошибка:", err))
            .finally(() => setLoading(false));
    }, []);

    return (
        <div className="space-y-6">
            <h1 className="text-3xl font-black text-slate-800 uppercase tracking-tighter italic">Команда</h1>

            <div className="bg-white rounded-[40px] border border-slate-200 overflow-hidden shadow-sm">
                <table className="w-full text-left">
                    <thead className="bg-slate-50 border-b border-slate-100">
                    <tr>
                        <th className="px-8 py-5 text-[10px] font-black text-slate-400 uppercase tracking-widest">Имя</th>
                        <th className="px-8 py-5 text-[10px] font-black text-slate-400 uppercase tracking-widest">Email</th>
                        <th className="px-8 py-5 text-[10px] font-black text-slate-400 uppercase tracking-widest text-right">Статус</th>
                    </tr>
                    </thead>
                    <tbody className="divide-y divide-slate-50">
                    {students.map((s) => (
                        <tr key={s.id} className="hover:bg-slate-50/50 transition-colors">
                            <td className="px-8 py-5 font-bold text-slate-700">{s.first_name} {s.last_name}</td>
                            <td className="px-8 py-5 text-slate-400 font-medium text-sm">{s.email}</td>
                            <td className="px-8 py-5 text-right">
                  <span className="bg-green-50 text-green-500 px-3 py-1 rounded-lg text-[10px] font-black uppercase">
                    {s.role}
                  </span>
                            </td>
                        </tr>
                    ))}
                    </tbody>
                </table>
            </div>
        </div>
    );
};

export default Students;