import { DuplicateIcon, LinkIcon } from "@heroicons/react/solid"

type ResultProps = {
  shortUrl: string
}

const Result = ({ shortUrl }: ResultProps) => (
  <div>
    <div className="relative">
      <LinkIcon className="h-5 w-5 absolute left-2.5 top-1/2 transform -translate-y-1/2 text-gray-400" />
      <input
        type="text"
        className="border-gray-300 border rounded w-full min-h-30 py-2 pr-2.5 pl-9 focus:ring-2 focus:ring-primary-light focus:outline-none focus:border-primary-light transition-shadow"
        placeholder="Enter a URL here"
        value={shortUrl.replace(/^https?:\/\//, "")}
        readOnly
      />
    </div>

    <button
      onClick={() => navigator.clipboard.writeText(shortUrl)}
      className="bg-primary text-white rounded font-semibold flex justify-center relative items-center w-full py-3 mt-3 disabled:opacity-50 focus:bg-primary-dark transition-colors"
    >
      <div className="flex items-center">
        <DuplicateIcon className="h-5 w-5 mr-2" />
        Copy
      </div>
    </button>
  </div>
)

export default Result
