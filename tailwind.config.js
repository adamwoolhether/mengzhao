/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./view/**/*.templ}", "./**/*.tmp"],
  safelist: [],
  plugins: [require("daisyui")],
  daisyui: {
    themes: ["synthwave"]
  }
}

