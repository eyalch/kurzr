const IndexPage = () => (
  <div className="text-center p-6">
    <h1 className="text-primary text-5xl font-bold mb-7">shrtr</h1>

    <form>
      <div className="relative">
        <svg
          className="h-5 w-5 absolute left-2.5 top-1/2 transform -translate-y-1/2 text-gray-400"
          fill="none"
          viewBox="0 0 24 24"
          stroke="currentColor"
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            strokeWidth={2}
            d="M13.828 10.172a4 4 0 00-5.656 0l-4 4a4 4 0 105.656 5.656l1.102-1.101m-.758-4.899a4 4 0 005.656 0l4-4a4 4 0 00-5.656-5.656l-1.1 1.1"
          />
        </svg>

        <input
          className="border-gray-300 border rounded w-full py-2 pr-2.5 pl-9 focus:ring-2 focus:ring-primary-light focus:outline-none focus:border-primary-light transition-shadow"
          type="text"
          placeholder="Enter a URL here"
        />
      </div>

      <button className="bg-primary text-white rounded font-semibold w-full py-3 mt-3">
        Shorten URL
      </button>
    </form>
  </div>
)

export default IndexPage
