# Bookings and Reservations

This is the repo for my bookings and reservations project

## Dependencies

- go 1.24.x
- [chi router](https://github.com/go-chi/chi)
- [session management (scs)](https://github.com/alexedwards/scs)
- Uses [nosurf](https://github.com/justinas/nosurf) for CSRF management
- template engine using [templ](https://github.com/a-h/templ)
- Tailwindcss v4

## Setup and Run

1. Build the project

   ```bash
   go mod tidy
   ```

2. Run the project

   ```bash
   # start tailwindcss compiling and watch process
   npm run watch:css

   # start the go project
   make run-local-server
   ```

3. Open the app in a browser

   Go to you browser and open the app on <http://localhost:7331>
