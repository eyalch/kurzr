# kurzr

kurzr is a simple URL shortener.

## Features

The main feature is to take a full (long) URL
and produce a short URL. Requesting the resulting short URL will return an HTTP
redirect status to the original full URL.

### Alias

When generating a short URL, you may optionally provide an alias which will be
used as the short identifier.

## Architecture

### API

The API is built with Go and implements the
[Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html).

### Web

The web UI is built with [React](https://reactjs.org/) and
[Next.js](https://nextjs.org/).

## Infrastructure

### Runtime

The project is running as a serverless function using
[Netlify Functions](https://www.netlify.com/products/functions/), which uses
[AWS Lambda](https://aws.amazon.com/lambda/) behind the scenes.

### Persistence

URLs are stored in [Redis](https://redis.io/).

### Security

#### reCAPTCHA

[reCAPTCHA](https://www.google.com/recaptcha/about/) is used to make it harder
for bots to use the site.

#### Rate Limiting

Redis is also utilized for rate limiting.
