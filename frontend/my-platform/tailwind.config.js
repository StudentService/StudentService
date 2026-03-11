/** @type {import('tailwindcss').Config} */
export default {
    content: ["./index.html", "./src/**/*.{js,jsx}"],
    theme: {
        extend: {
            colors: {
                brand: {
                    DEFAULT: '#15808d', // Тот самый бирюзовый из «Капитанов»
                    dark: '#0e5a63',
                    light: '#e8f3f4',
                }
            }
        },
    },
    plugins: [],
}