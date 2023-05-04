/** @type {import('tailwindcss').Config} */
export default {
    content: ['./src/**/*.{html,js,svelte,ts}'],
    theme: {
        extend: {},
    },

    daisyui: {
        themes: [
            'light',
            'dark',
            'luxury',
            // Luxury theme with cyberpunk style
            {'cyberlux':{
                "color-scheme": "dark",
                fontFamily: "ui-monospace,SFMono-Regular,Menlo,Monaco,Consolas,Liberation Mono,Courier New,monospace",
                primary: "#ffffff",
                secondary: "#152747",
                accent: "#513448",
                neutral: "#171618",
                "neutral-content": "#dca54c",
                "base-100": "#09090b",
                "base-200": "#171618",
                "base-300": "#2e2d2f",
                "base-content": "#dca54c",
                info: "#66c6ff",
                success: "#87d039",
                warning: "#e2d562",
                error: "#ff6f6f",
                "--rounded-box": "0",
                "--rounded-btn": "0",
                "--rounded-badge": "0",
                "--tab-radius": "0",
            }},
            // Dark theme with cyberpunk style
            {'cyberdark':{
                "color-scheme": "dark",
                fontFamily: "ui-monospace,SFMono-Regular,Menlo,Monaco,Consolas,Liberation Mono,Courier New,monospace",
                primary: "#661AE6",
                "primary-content": "#ffffff",
                secondary: "#D926AA",
                "secondary-content": "#ffffff",
                accent: "#1FB2A5",
                "accent-content": "#ffffff",
                neutral: "#191D24",
                "neutral-focus": "#111318",
                "neutral-content": "#A6ADBB",
                "base-100": "#2A303C",
                "base-200": "#242933",
                "base-300": "#20252E",
                "base-content": "#A6ADBB",
                "--rounded-box": "0",
                "--rounded-btn": "0",
                "--rounded-badge": "0",
                "--tab-radius": "0",
            }},
        ]
    },
    plugins: [require("daisyui")],
}

