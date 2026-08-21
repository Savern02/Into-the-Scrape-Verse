import type { Config } from "tailwindcss";

export default {
  darkMode: ["class"],
  content: [
    "./pages/**/*.{js,ts,jsx,tsx,mdx}",
    "./components/**/*.{js,ts,jsx,tsx,mdx}",
    "./app/**/*.{js,ts,jsx,tsx,mdx}",
    "./src/**/*.{js,ts,jsx,tsx,mdx}",
  ],
  theme: {
    extend: {
      colors: {
        seaweed: {
        50: "#eef7f3",
        100: "#dcefe6",
        200: "#b9dfcd",
        300: "#96cfb4",
        400: "#73bf9c",
        500: "#50af83",
        600: "#408c69",
        700: "#30694e",
        800: "#204634",
        900: "#10231a",
        950: "#0b1812",
      },

      charcoal: {
        50: "#f3f2f2",
        100: "#e7e5e4",
        200: "#cecbca",
        300: "#b6b0af",
        400: "#9d9695",
        500: "#857c7a",
        600: "#6a6362",
        700: "#504a49",
        800: "#353231",
        900: "#1b1918",
        950: "#131111",
      },

      "light-yellow": {
        50: "#fffce5",
        100: "#fff9cc",
        200: "#fff399",
        300: "#ffed66",
        400: "#ffe733",
        500: "#ffe100",
        600: "#ccb400",
        700: "#998700",
        800: "#665a00",
        900: "#332d00",
        950: "#242000",
      },

      "banana-cream": {
        50: "#fffce5",
        100: "#fff9cc",
        200: "#fff399",
        300: "#ffed66",
        400: "#ffe733",
        500: "#ffe100",
        600: "#ccb400",
        700: "#998700",
        800: "#665a00",
        900: "#332d00",
        950: "#242000",
      },

      tomato: {
        50: "#ffe9e5",
        100: "#ffd4cc",
        200: "#ffa899",
        300: "#ff7d66",
        400: "#ff5233",
        500: "#ff2600",
        600: "#cc1f00",
        700: "#991700",
        800: "#660f00",
        900: "#330800",
        950: "#240500",
      },
        background: "hsl(var(--background))",
        foreground: "hsl(var(--foreground))",
        card: {
          DEFAULT: "hsl(var(--card))",
          foreground: "hsl(var(--card-foreground))",
        },
        popover: {
          DEFAULT: "hsl(var(--popover))",
          foreground: "hsl(var(--popover-foreground))",
        },
        primary: {
          DEFAULT: "hsl(var(--primary))",
          foreground: "hsl(var(--primary-foreground))",
        },
        secondary: {
          DEFAULT: "hsl(var(--secondary))",
          foreground: "hsl(var(--secondary-foreground))",
        },
        muted: {
          DEFAULT: "hsl(var(--muted))",
          foreground: "hsl(var(--muted-foreground))",
        },
        accent: {
          DEFAULT: "hsl(var(--accent))",
          foreground: "hsl(var(--accent-foreground))",
        },
        destructive: {
          DEFAULT: "hsl(var(--destructive))",
          foreground: "hsl(var(--destructive-foreground))",
        },
        border: "hsl(var(--border))",
        input: "hsl(var(--input))",
        ring: "hsl(var(--ring))",
        chart: {
          "1": "hsl(var(--chart-1))",
          "2": "hsl(var(--chart-2))",
          "3": "hsl(var(--chart-3))",
          "4": "hsl(var(--chart-4))",
          "5": "hsl(var(--chart-5))",
        },
      },
      borderRadius: {
        lg: "var(--radius)",
        md: "calc(var(--radius) - 2px)",
        sm: "calc(var(--radius) - 4px)",
      },
    },
  },
  plugins: [require("tailwindcss-animate")],
} satisfies Config;
