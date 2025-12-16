# tmdb - The movie db CLI tool

CLI tool to discover hidden gems in your streaming libraries. Fetches data from TMDB and your streaming services to 
help you find movies and shows you might have missed. 

## Example

```bash
tmdb actor "Leonardo DiCaprio"
```
![tmdb-actor-example](docs/example.png)

## Build and Install

```bash
go build -o tmdb cmd/main.go # or make build

./tmdb
```

## Setup

1. Get an API key from TMDB: https://www.themoviedb.org/settings/api
2. Copy `.env.example` to `.env`
3. Add your TMDB API key and adjust your preferences:

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

# list available genres in your language/region
./tmdb genres

# show popular movies
./tmdb popular --genre action

# Show filmographies of a specific actor
./tmdb actor "Nicolas Cage"

# List popular actors
./tmdb actor

# Search for actors with partial name match and show all options
./tmdb actor "tom"

# Select specific actor from multiple results by index
./tmdb actor "Megan Fox" 1

# Search for movies with "star" in the title like Star Wars, Star Trek, etc.
./tmdb search star

# Show top rated shows
./tmdb shows --min-rating 8.0
```

## Actor Command - Multiple Results

When searching for an actor that returns multiple results (e.g., "Tom" returns multiple actors named Tom):

1. **Without index**: Shows a list of matching actors
   ```bash
   tmdb actor "tom"
   # Output shows: [1] Tom Hanks, [2] Tom Hardy, [3] Tom Cruise, etc.
   ```

2. **With index**: Directly selects and fetches filmography for a specific actor
   ```bash
   tmdb actor "tom" 1        # Fetch Tom Hanks filmography
   tmdb actor "tom" 3        # Fetch Tom Cruise filmography
   ```

The index corresponds to the order shown in the search results (1-based indexing).

## Missing Features

- [ ] caching of API responses to reduce load times and API calls
- [ ] "fomo" command. List movies leaving streaming services soon.
- [ ] generate .env file with setup command (interactive shell)
- [ ] list possible streaming services