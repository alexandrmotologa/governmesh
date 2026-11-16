/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        mesh: {
          bg: "#090d16",
          card: "#0f172a",
          border: "#1e293b",
          primary: "#3b82f6",
          accent: "#06b6d4",
          success: "#10b981",
          warning: "#f59e0b",
          danger: "#ef4444",
        }
      }
    },
  },
  plugins: [],
}
