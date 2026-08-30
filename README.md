# reddit-pwa-client

Mobile-friendly reddit client: home feed, threaded comments, voting. Personal use only, not for distribution.

Go stdlib backend + htmx/server-rendered HTML frontend. Installable as a PWA.

## Setup

1. Copy `.env.example` to `.env`.
2. Log into old.reddit.com in a browser, DevTools → Network → any request → copy `reddit_session` cookie value.
3. Fill in `.env`:
   ```env
   REDDIT_USERNAME=your_username
   REDDIT_SESSION_COOKIE=<copied value>
   ```

## Run

```bash
go run main.go
```

Hot reload during development:

```bash
go install github.com/air-verse/air@latest
air
```

Visit `http://localhost:8080/home`.

## Notes

- Session cookie expires eventually — repeat step 2-3 above and restart when requests start failing.
- No offline support beyond installability.
