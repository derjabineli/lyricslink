module.exports = {
  content: [
    "./views/**/*.html", // Ensure Tailwind scans all HTML files
    "./static/**/*.js", // Ensure Tailwind scans all JavaScript files
  ],
  theme: {
    extend: {
      fontFamily: {
        sans: ["var(--font-sans)"],
      },
    },
  },
  plugins: [],
};
