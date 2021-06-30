import { useState } from "react"
import Form from "../components/Form"
import Result from "../components/Result"

const IndexPage = () => {
  const [url, setUrl] = useState("")
  const [shortUrl, setShortUrl] = useState<string>()

  return (
    <div className="p-6">
      <h1 className="text-primary text-5xl text-center font-bold mb-7">
        shrtr
      </h1>

      {shortUrl ? (
        <Result shortUrl={shortUrl} />
      ) : (
        <Form url={url} setUrl={setUrl} setShortUrl={setShortUrl} />
      )}
    </div>
  )
}

export default IndexPage
