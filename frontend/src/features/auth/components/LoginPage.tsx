// src/features/auth/components/AuthPage.tsx
import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import api from '@/shared/api'
import { useAuthStore } from '@/features/auth/store'

import {
    Box,
    Button,
    Container,
    TextField,
    Typography,
    Alert,
    Paper,
    CircularProgress,
    InputAdornment,
    IconButton,
    Link as MuiLink,
    Divider,
} from '@mui/material'
import { Visibility, VisibilityOff } from '@mui/icons-material'

export default function AuthPage() {
    const [mode, setMode] = useState<'login' | 'register'>('login')

    // Общие поля
    const [login, setLogin] = useState('')
    const [password, setPassword] = useState('')
    const [showPassword, setShowPassword] = useState(false)

    // Поля только для регистрации (из твоего профиля)
    const [firstName, setFirstName] = useState('')
    const [lastName, setLastName] = useState('')
    const [group, setGroup] = useState('')
    const [segment, setSegment] = useState('')
    const [semester, setSemester] = useState('')

    const [error, setError] = useState('')
    const [loading, setLoading] = useState(false)

    const navigate = useNavigate()

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault()
        setError('')
        setLoading(true)

        try {
            if (mode === 'login') {
                // Вход
                const res = await api.post('/auth/login', {
                    login: login.trim(),
                    password,
                })

                const { token, user } = res.data
                useAuthStore.getState().login(token, user)
                navigate('/dashboard', { replace: true })
            } else {
                // Регистрация
                const payload = {
                    firstName: firstName.trim(),
                    lastName: lastName.trim(),
                    group: group.trim(),
                    segment: segment.trim(),
                    semester: semester.trim(),
                    login: login.trim(),
                    password,
                }

                // Замени путь на реальный, если он другой
                const res = await api.post('/auth/register', payload) // или /students, /signup и т.д.

                // Если бэк возвращает токен сразу — логиним
                if (res.data.token) {
                    useAuthStore.getState().login(res.data.token, res.data.user)
                    navigate('/dashboard', { replace: true })
                } else {
                    // Если нужно отдельно логиниться
                    alert('Регистрация успешна! Теперь войдите.')
                    setMode('login')
                }

                setError('')
            }
        } catch (err: unknown) {
            const message = err instanceof Error ? err.message : 'Ошибка. Попробуйте позже'
            setError(message)
        } finally {
            setLoading(false)
        }
    }

    const toggleMode = () => {
        setMode(mode === 'login' ? 'register' : 'login')
        setError('')
    }

    return (
        <Container maxWidth="sm" sx={{ mt: 8, mb: 8 }}>
            <Paper elevation={6} sx={{ p: 5, borderRadius: 3 }}>
                <Typography variant="h4" align="center" gutterBottom>
                    {mode === 'login' ? 'Вход' : 'Регистрация'}
                </Typography>

                {error && <Alert severity="error" sx={{ mb: 3 }}>{error}</Alert>}

                <Box component="form" onSubmit={handleSubmit} noValidate>
                    {mode === 'register' && (
                        <>
                            <TextField
                                margin="normal"
                                required
                                fullWidth
                                label="Имя"
                                value={firstName}
                                onChange={e => setFirstName(e.target.value)}
                                disabled={loading}
                            />
                            <TextField
                                margin="normal"
                                required
                                fullWidth
                                label="Фамилия"
                                value={lastName}
                                onChange={e => setLastName(e.target.value)}
                                disabled={loading}
                            />
                            <TextField
                                margin="normal"
                                required
                                fullWidth
                                label="Группа"
                                value={group}
                                onChange={e => setGroup(e.target.value)}
                                disabled={loading}
                            />
                            <TextField
                                margin="normal"
                                required
                                fullWidth
                                label="Сегмент"
                                value={segment}
                                onChange={e => setSegment(e.target.value)}
                                disabled={loading}
                            />
                            <TextField
                                margin="normal"
                                required
                                fullWidth
                                label="Семестр"
                                value={semester}
                                onChange={e => setSemester(e.target.value)}
                                disabled={loading}
                            />
                            <Divider sx={{ my: 3 }} />
                        </>
                    )}

                    <TextField
                        margin="normal"
                        required
                        fullWidth
                        label="Логин"
                        autoFocus
                        value={login}
                        onChange={e => setLogin(e.target.value.trim())}
                        disabled={loading}
                    />

                    <TextField
                        margin="normal"
                        required
                        fullWidth
                        label="Пароль"
                        type={showPassword ? 'text' : 'password'}
                        value={password}
                        onChange={e => setPassword(e.target.value)}
                        disabled={loading}
                        InputProps={{
                            endAdornment: (
                                <InputAdornment position="end">
                                    <IconButton
                                        onClick={() => setShowPassword(!showPassword)}
                                        edge="end"
                                        disabled={loading}
                                    >
                                        {showPassword ? <VisibilityOff /> : <Visibility />}
                                    </IconButton>
                                </InputAdornment>
                            ),
                        }}
                    />

                    <Button
                        type="submit"
                        fullWidth
                        variant="contained"
                        size="large"
                        disabled={loading}
                        sx={{ mt: 4, py: 1.5 }}
                    >
                        {loading ? (
                            <>
                                <CircularProgress size={24} sx={{ mr: 1 }} />
                                {mode === 'login' ? 'Вход...' : 'Регистрация...'}
                            </>
                        ) : mode === 'login' ? 'Войти' : 'Зарегистрироваться'}
                    </Button>

                    <Box textAlign="center" sx={{ mt: 3 }}>
                        <MuiLink
                            component="button"
                            variant="body2"
                            onClick={toggleMode}
                            sx={{ cursor: 'pointer' }}
                        >
                            {mode === 'login'
                                ? 'Нет аккаунта? Зарегистрироваться'
                                : 'Уже есть аккаунт? Войти'}
                        </MuiLink>
                    </Box>
                </Box>
            </Paper>
        </Container>
    )
}