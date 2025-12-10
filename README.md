# tmdb - the movie db CLI tool

CLI tool to discover hidden gems in your streaming libraries. Fetches data from TMDB and your streaming services to 
help you find movies and shows you might have missed. 

## Build and Install

```bash
go build -o tmdb cmd/main.go # or make build

./tmdb
```

## Setup

1. Copy `.env.example` to `.env`
2. Add your TMDB API key and adjust your preferences:

    ```
    # .env file
    TMDB_API_KEY=your_api_key_here
    PROVIDERS=Netflix,DisneyPlus,Wow,RtlPlus
    REGION=DE
    MIN_RATING=7.5
    MIN_VOTES=1000
    API_TIMEOUT_SECONDS=20
    ```

## Usage

```bash
# Show top rated movies filtered by your .env settings or cli options
./tmdb top --min-rating 6.0

# Show filmographies of a specific actor
./tmdb actor "Nicolas Cage"

# Show top rated shows
./tmdb shows --min-rating 8.0
```

## Missing Features

- [ ] fomo command. List movies leaving streaming services soon.
- [ ] add movie search by title
- [ ] add more filters like genre, release year, etc.
- [ ] caching of API responses to reduce load times and API calls
- [ ] add tests
