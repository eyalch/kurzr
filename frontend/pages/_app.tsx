import type { AppProps } from "next/app"
import "tailwindcss/tailwind.css"

export default function MyApp({ Component, pageProps }: AppProps) {
  return (
    <main className="max-w-xl mx-auto">
      <Component {...pageProps} />
    </main>
  )
}
