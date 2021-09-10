import { faLink, faTag } from "@fortawesome/free-solid-svg-icons"
import axios from "axios"
import React, { useEffect, useState } from "react"
import { SubmitHandler, useForm } from "react-hook-form"
import Button from "./Button"
import Input from "./Input"

const http = axios.create({
  baseURL: process.env.NEXT_PUBLIC_API_URL ?? "",
})

type FormProps = {
  onSuccess: (url: string, shortUrl: string) => void
}

type IFormInput = {
  url: string
  alias: string
}

// TODO: Add share option (on mobile)
// TODO: Add QR code
// TODO: Add unit tests

const recaptchaKey = process.env.NEXT_PUBLIC_RECAPTCHA_KEY!

const isValidUrl = (string: string) => {
  try {
    new URL(string)
    return true
  } catch {
    return false
  }
}

const loadScriptByURL = (url: string, callback: () => void) => {
  const script = document.createElement("script")
  script.type = "text/javascript"
  script.src = url
  if (callback) script.onload = callback

  document.body.appendChild(script)
}

const Form = ({ onSuccess }: FormProps) => {
  const [loading, setLoading] = useState(true)

  const {
    register,
    handleSubmit,
    watch,
    setError,
    formState: { errors },
  } = useForm<IFormInput>({ mode: "all" })

  const url = watch("url", "")
  const urlWithSchema = url.includes("://") ? url : `https://${url}`

  const shortenUrl: SubmitHandler<IFormInput> = async ({ alias, url }) => {
    setLoading(true)

    const token = await window.grecaptcha.execute(recaptchaKey, {
      action: "submit",
    })

    try {
      const { data } = await http.post("/api", {
        url: urlWithSchema,
        alias,
        token,
      })
      onSuccess(url, data.short_url)
    } catch (error) {
      setLoading(false)
      if (error.response.data.code === "ERR_DUPLICATE_KEY") {
        setError("alias", {
          type: "duplicate",
          message: "Alias is not available",
        })
      }
    }
  }

  useEffect(() => {
    loadScriptByURL(
      `https://www.google.com/recaptcha/api.js?render=${recaptchaKey}`,
      () => setLoading(false)
    )
  }, [])

  return (
    <form onSubmit={handleSubmit(shortenUrl)} autoComplete="off">
      <Input
        icon={faLink}
        placeholder="Enter a URL here"
        {...register("url")}
        className="mb-4"
      />

      <Input
        icon={faTag}
        placeholder="Alias (optional)"
        {...register("alias", {
          pattern: {
            value: /^[A-Za-z0-9]+$/,
            message: "Use only letters and numbers",
          },
        })}
        error={errors.alias?.message}
        className="mb-4"
      />

      <Button
        label="Shorten URL"
        disabled={!isValidUrl(urlWithSchema)}
        loading={loading}
      />
    </form>
  )
}

export default Form
