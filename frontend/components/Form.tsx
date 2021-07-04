import { faLink } from "@fortawesome/free-solid-svg-icons"
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome"
import { FormEvent, useState } from "react"
import Button from "./Button"

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
        <FontAwesomeIcon
          icon={faLink}
          className="h-5 w-5 absolute left-2.5 top-1/2 transform -translate-y-1/2 text-gray-400"
        />

        <input
          type="text"
          className="border-gray-300 border rounded w-full min-h-30 py-2 pr-2.5 pl-9 focus:ring-2 focus:ring-primary-light focus:outline-none focus:border-primary-light transition-shadow"
          placeholder="Enter a URL here"
          value={url}
          onChange={(event) => setUrl(event.target.value)}
        />
      </div>

      <div className="flex mt-3">
        <Button
          label="Shorten URL"
          disabled={!isValidUrl(urlWithSchema)}
          loading={loading}
        ></Button>
      </div>
    </form>
  )
}

export default Form
