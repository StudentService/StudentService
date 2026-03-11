import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import LoginPage from './pages/LoginPage';
import RegisterPage from './pages/RegisterPage';
import MainLayout from './components/MainLayout';
import Dashboard from './pages/Dashboard';
import Challenges from './pages/Challenges';

// Компонент-защитник: не дает войти в систему без токена
const ProtectedRoute = ({ children }) => {
    const token = localStorage.getItem('access_token');
    if (!token) {
        return <Navigate to="/login" replace />;
    }
    return children;
};

// Компонент-инвертор: не дает зайти на логин/регистрацию, если уже авторизован
const PublicRoute = ({ children }) => {
    const token = localStorage.getItem('access_token');
    if (token) {
        return <Navigate to="/dashboard" replace />;
    }
    return children;
};

function App() {
    return (
        <BrowserRouter>
            <Routes>
                {/* Публичные страницы: логин и регистрация */}
                <Route path="/login" element={
                    <PublicRoute>
                        <LoginPage />
                    </PublicRoute>
                } />
                <Route path="/register" element={
                    <PublicRoute>
                        <RegisterPage />
                    </PublicRoute>
                } />

                {/* Защищенная часть приложения */}
                <Route
                    element={
                        <ProtectedRoute>
                            <MainLayout />
                        </ProtectedRoute>
                    }
                >
                    <Route path="/dashboard" element={<Dashboard />} />
                    <Route path="/students" element={<div className="p-8 text-xl">Страница студентов</div>} />
                    <Route path="/metrics" element={<div className="p-8 text-xl">Конструктор метрик</div>} />
                    <Route path="/challenges" element={<Challenges />} />
                    {/* Редирект по умолчанию на дашборд */}
                    <Route path="/" element={<Navigate to="/dashboard" replace />} />
                </Route>

                {/* Все неизвестные пути — на логин */}
                <Route path="*" element={<Navigate to="/login" replace />} />
            </Routes>
        </BrowserRouter>
    );
}

export default App;