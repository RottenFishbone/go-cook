/** @type {import('tailwindcss').Config} */
export default {
    content: ['./src/**/*.{html,js,svelte,ts}'],
    theme: {
        extend: {},
    },

    daisyui: {
        themes: [
			{'ultracupcake':{
                ...require("daisyui/src/colors/themes")["[data-theme=cupcake]"],
                fontFamily: "ui-monospace,SFMono-Regular,Menlo,Monaco,Consolas,Liberation Mono,Courier New,monospace",
				"base-100": "#ffffff",
				"base-200": "#f6f8fa",
				"base-300": "#eaeef2",
				"primary":"#661AE6",
				"accent":"#D926AA",
				"info":"#0284c7",
				"--rounded-box": "0",
                "--rounded-btn": "0",
                "--rounded-badge": "0",
                "--tab-radius": "0",
			}},
            // Lightented forest theme with cyberpunk styling
            {'ultraforest':{
                ...require("daisyui/src/colors/themes")["[data-theme=forest]"],
                fontFamily: "ui-monospace,SFMono-Regular,Menlo,Monaco,Consolas,Liberation Mono,Courier New,monospace",
                "base-100": "#2c2727",
				"primary-content" : "#093611",
				"--rounded-box": "0",
                "--rounded-btn": "0",
                "--rounded-badge": "0",
                "--tab-radius": "0",
            }},
        ],
    },
    plugins: [require("daisyui")],
}

