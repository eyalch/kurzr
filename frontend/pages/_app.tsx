import type { AppProps } from "next/app"
import "tailwindcss/tailwind.css"

function MyApp({ Component, pageProps }: AppProps) {
  return (
    <main className="max-w-xl mx-auto">
      <Component {...pageProps} />
    </main>
  )
}

export default MyApp
