module.exports = {
  purge: ["./pages/**/*.{js,ts,jsx,tsx}"],
  darkMode: false, // or 'media' or 'class'
  theme: {
    extend: {
      colors: {
        primary: {
          light: "#89C2D9",
          DEFAULT: "#3181AF",
          dark: "#1F526F",
        },
      },
      scale: {
        "-1": "-1",
      },
    },
    borderRadius: {
      DEFAULT: "0.375rem",
    },
  },
  variants: {
    extend: {
      opacity: ["disabled"],
    },
  },
  plugins: [],
}
