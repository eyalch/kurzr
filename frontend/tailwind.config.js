module.exports = {
  mode: "jit",
  purge: ["./pages/**/*.{js,ts,jsx,tsx}"],
  darkMode: false, // or 'media' or 'class'
  theme: {
    extend: {
      colors: {
        primary: {
          light: "#89C2D9",
          DEFAULT: "#2A6F97",
        },
      },
    },
    borderRadius: {
      DEFAULT: "0.375rem",
    },
  },
  variants: {
    extend: {},
  },
  plugins: [],
}
