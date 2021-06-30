import { LinkIcon } from "@heroicons/react/solid"
import { FormEvent, useState } from "react"

type FormProps = {
  url: string
  setUrl: (url: string) => void
  setShortUrl: (shortUrl: string) => void
}

const isValidUrl = (string: string) => {
  try {
    new URL(string)
    return true
  } catch {
    return false
  }
}

const Form = ({ url, setUrl, setShortUrl }: FormProps) => {
  const [loading, setLoading] = useState(false)

  const urlWithSchema = url.includes("://") ? url : `https://${url}`

  const shortenUrl = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault()

    setLoading(true)

    try {
      const response = await fetch("/api", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ url: urlWithSchema }),
      })
      const data = await response.json()

      setShortUrl(data.url)
    } catch {
      setLoading(false)
    }
  }

  return (
    <form onSubmit={shortenUrl}>
      <div className="relative">
        <LinkIcon className="h-5 w-5 absolute left-2.5 top-1/2 transform -translate-y-1/2 text-gray-400" />
        <input
          type="text"
          className="border-gray-300 border rounded w-full min-h-30 py-2 pr-2.5 pl-9 focus:ring-2 focus:ring-primary-light focus:outline-none focus:border-primary-light transition-shadow"
          placeholder="Enter a URL here"
          value={url}
          onChange={(event) => setUrl(event.target.value)}
        />
      </div>

      <div className="flex mt-3">
        <button
          disabled={!isValidUrl(urlWithSchema)}
          className="bg-primary text-white rounded font-semibold flex justify-center relative items-center w-full py-3 disabled:opacity-50 focus:bg-primary-dark transition-colors"
        >
          {loading && (
            <svg
              className="absolute inset-x-0 w-full animate-spin h-5 text-white"
              fill="none"
              viewBox="0 0 24 24"
            >
              <circle
                className="opacity-25"
                cx="12"
                cy="12"
                r="10"
                stroke="currentColor"
                strokeWidth="4"
              />
              <path
                className="opacity-75"
                fill="currentColor"
                d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
              />
            </svg>
          )}

          <div className={`flex items-center ${loading ? "opacity-0" : ""}`}>
            Shorten URL
          </div>
        </button>
      </div>
    </form>
  )
}

export default Form
