import {
  faCopy,
  faExternalLinkAlt,
  faLink,
} from "@fortawesome/free-solid-svg-icons"
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome"
import React from "react"
import Button from "./Button"

type ResultProps = {
  shortUrl: string
}

const Result = ({ shortUrl }: ResultProps) => (
  <div>
    <div className="relative">
      <FontAwesomeIcon
        icon={faLink}
        className="h-5 w-5 absolute left-2.5 top-1/2 transform -translate-y-1/2 text-gray-400"
      />
      <input
        type="text"
        className="border-gray-300 border rounded w-full min-h-30 py-2 pr-2.5 pl-9 focus:ring-2 focus:ring-primary-light focus:outline-none focus:border-primary-light transition-shadow"
        placeholder="Enter a URL here"
        value={shortUrl.replace(/^https?:\/\//, "")}
        readOnly
      />
    </div>

    <div className="mt-3 grid grid-cols-2 gap-3">
      <Button
        icon={faCopy}
        label="Copy"
        className=""
        onClick={() => navigator.clipboard.writeText(shortUrl)}
      />

      <Button
        icon={faExternalLinkAlt}
        label="Visit"
        href={shortUrl}
        target="_blank"
        className=""
      />
    </div>
  </div>
)

export default Result
