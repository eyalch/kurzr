module.exports = {
  async rewrites() {
    return [
      {
        source: "/api/:path*",
        destination: "http://localhost:4000/:path*", // Proxy to Backend
      },
    ]
  },
}
