/** @type {import('tailwindcss').Config} */
export default {
    content: ['./src/**/*.{html,js,svelte,ts}'],
    theme: {
        extend: {},
    },

    daisyui: {
        themes: [
            'cupcake',
            // Lightented forest theme with cyberpunk styling
            {'ultraforest':{
                ...require("daisyui/src/colors/themes")["[data-theme=forest]"],
                fontFamily: "ui-monospace,SFMono-Regular,Menlo,Monaco,Consolas,Liberation Mono,Courier New,monospace",
                "base-100": "#2c2727",
                "--rounded-box": "0",
                "--rounded-btn": "0",
                "--rounded-badge": "0",
                "--tab-radius": "0",
            }},
        ],
    },
    plugins: [require("daisyui")],
}

