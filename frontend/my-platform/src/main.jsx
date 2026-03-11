import React from 'react'
import ReactDOM from 'react-dom/client'
import './index.css' // Сначала стили!
import App from './App' // Потом приложение

ReactDOM.createRoot(document.getElementById('root')).render(
    <React.StrictMode>
        <App />
    </React.StrictMode>,
)
