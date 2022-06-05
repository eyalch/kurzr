import Head from "next/head"
import React, { PointerEvent, useState } from "react"
import Form from "../components/Form"
import Result from "../components/Result"

export default function IndexPage() {
  const [, setUrl] = useState("")
  const [shortUrl, setShortUrl] = useState<string>()

  const resetUrl = (event: PointerEvent<HTMLAnchorElement>) => {
    event.preventDefault()

    setShortUrl("")
    setUrl("")
  }

  const onSuccess = (url: string, shortUrl: string) => {
    setUrl(url)
    setShortUrl(shortUrl)
  }

  return (
    <>
      <Head>
        <title>kurzr</title>
      </Head>

      <div className="p-3">
        <h1 className="text-primary text-5xl text-center font-bold my-7">
          <a onClick={resetUrl} href="/">
            kurzr
          </a>
        </h1>

        {shortUrl ? (
          <Result shortUrl={shortUrl} />
        ) : (
          <Form onSuccess={onSuccess} />
        )}
      </div>
    </>
  )
}
