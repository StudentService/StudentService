import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import MainLayout from './components/MainLayout';
import Dashboard from './pages/Dashboard';
import ProfilePage from './pages/ProfilePage';
import QuestionnairePage from './pages/QuestionnairePage'; // Импорт
import ChallengesPage from './pages/ChallengesPage'; // Импорт
import LoginPage from './pages/LoginPage';
import CalendarPage from "./pages/CalendarPage.jsx";
import GradesPage from "./pages/GradesPage.jsx";
import ActivitiesPage from "./pages/ActivitiesPage.jsx";
function App() {
    return (
        <Router>
            <Routes>
                <Route path="/login" element={<LoginPage />} />

                <Route path="/" element={<MainLayout />}>
                    <Route index element={<Navigate to="/dashboard" />} />
                    <Route path="dashboard" element={<Dashboard />} />
                    <Route path="profile" element={<ProfilePage />} />
                    <Route path="form" element={<QuestionnairePage />} />
                    <Route path="challenges" element={<ChallengesPage />} />
                    <Route path="calendar" element={<CalendarPage />} />
                    <Route path="grades" element={<GradesPage />} />
                    <Route path="activities" element={<ActivitiesPage />} />

                </Route>
            </Routes>
        </Router>
    );
}

export default App;