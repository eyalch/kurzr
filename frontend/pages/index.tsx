import React, { PointerEvent, useState } from "react"
import Form from "../components/Form"
import Result from "../components/Result"

const IndexPage = () => {
  const [url, setUrl] = useState("")
  const [shortUrl, setShortUrl] = useState<string>()

  const resetUrl = (event: PointerEvent<HTMLAnchorElement>) => {
    event.preventDefault()

    setShortUrl("")
    setUrl("")
  }

  return (
    <div className="p-3">
      <a onClick={resetUrl} href="/">
        <h1 className="text-primary text-5xl text-center font-bold my-7">
          shrtr
        </h1>
      </a>

      {shortUrl ? (
        <Result shortUrl={shortUrl} />
      ) : (
        <Form url={url} setUrl={setUrl} setShortUrl={setShortUrl} />
      )}
    </div>
  )
}

export default IndexPage
