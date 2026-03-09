import { AppBar, Toolbar, Typography, Button, Box } from '@mui/material'
import { Link as RouterLink, useNavigate } from 'react-router-dom'
import { useAuthStore } from '@/features/auth/store'
import { Outlet } from 'react-router-dom'

export default function MainLayout() {
    const { logout } = useAuthStore()
    const navigate = useNavigate()

    return (
        <>
            <AppBar position="fixed">
                <Toolbar>
                    <Typography variant="h6" sx={{ flexGrow: 1 }}>
                        Платформа
                    </Typography>
                    <Button color="inherit" component={RouterLink} to="/profile">
                        Мой профиль
                    </Button>
                    <Button color="inherit" component={RouterLink} to="/dashboard">
                        Личный кабинет
                    </Button>
                    <Button color="inherit" component={RouterLink} to="/calendar">
                        Календарь
                    </Button>
                    <Button
                        color="inherit"
                        onClick={() => {
                            logout()
                            navigate('/login', { replace: true })
                        }}
                    >
                        Выход
                    </Button>
                </Toolbar>
            </AppBar>

            <Box sx={{ pt: 8, p: 3 }}>
                <Outlet />
            </Box>
        </>
    )
}