import { createBrowserRouter, Navigate } from 'react-router-dom'
import ProtectedRoute from '@/app/ProtectedRoute' // или твой путь

import LoginPage from '@/features/auth/components/LoginPage'
import ProfilePage from '@/features/profile/ProfilePage'     // ← правильный путь
import DashboardPage from '@/features/dashboard/DashboardPage' // ← правильный путь
import CalendarPage from '@/features/calendar/CalendarPage'   // ← правильный путь

// Layout с навбаром
import MainLayout from '@/shared/ui/layouts/MainLayout'

export const router = createBrowserRouter([
    // Публичный маршрут — вход (без навбара и без проверки)
    {
        path: '/login',
        element: <LoginPage />,
    },

    {
        element: <ProtectedRoute />,
        children: [
            {
                element: <MainLayout />, // ← здесь появляется навбар
                children: [
                    // Если зашли на корень сайта — редирект на дашборд
                    { index: true, element: <Navigate to="/dashboard" replace /> },

                    // Защищённые страницы
                    { path: 'dashboard', element: <DashboardPage /> },
                    { path: 'profile',   element: <ProfilePage />   },
                    { path: 'calendar',  element: <CalendarPage />  },
                ],
            },
        ],
    },

    // 404 на всё остальное
    {
        path: '*',
        element: <h1>404 — Страница не найдена</h1>,
    },
])