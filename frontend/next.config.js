module.exports = (phase, { defaultConfig }) => {
  if (process.env.NEXT_PUBLIC_API_URL) {
    return {
      ...defaultConfig,
      async rewrites() {
        return [
          {
            source: "/api/:path*",
            destination: `${process.env.NEXT_PUBLIC_API_URL}/:path*`,
          },
        ]
      },
    }
  }

  return defaultConfig
}
