# AQuickQuestion

- backend:
  - fiber
  - aargon2
  - jwt:
    - refresh token (httpOnly, httpsOnly, cookie ~30d)
      when the token expires the user must authenticate again
    - access token (saved in the local process ~15min)
      this token can be lost and remade on demand as long the refresh token is valid
      the server will first check if the access token is valid/present, if not,
      it will create a new access token that the client can use in the following calls
      (usually a dedicatated endpoint /api/token/refresh)
- db:
  - pg
- frontend:
  - templ
  - htmx
  - alpine
  - ts
