## Documentation

1. Chapter 1: Introduction to the book

2. Chapter 2: Getting Started/Project structure creation
   - 2.1: project structure setup
     1. Setup project structure
        - `bin`: contain compiled application binaries, ready for deployment to a production server
        - `cmd/api`: contain application-specific code i.e. code for running the server, reading and writing HTTP requests, and managing authentication
        - `internal`: contains _ancillary_ packages using by our API i.e. database interation code, data validation, sending emails etc. **Note: any code which isn't application-specific and can potentially be reused will live here. The code under `cmd/api` will import packages in the `internal` directory (but never the other way around). Any packages under this directory can only be imported by code inside the parent of the `internal` directory i.e. any code inside `greenlight`, else any package under `internal` cannot be imported by code outside of our project. This prevents other codebases from importing and relying on the potentially (unversioned and unsupported) packages in our internal directory, even if the project code is publicly available somewhere like Github.**
   - 2.2: Basic HTTP Server
     1. set up config and server structs
   - 2.3: API Endpoints and RESTful Routing
     1. installed and set [`httprouter`](github.com/julienschmidt/httprouter@v1.3.0) for routing

3. Chapter 3: Sending JSON (Javascript Object Notation) Responses
   - 3.1: Fixed-Format JSON
     1. updated healthcheckHandler to send JSON response
   - 3.2 JSON Encoding
     1. updated `healthcheckHandler` to use the `json.Marshal()` function to encode the JSON
     2. set a new function to `writeJSON` in `helpers.go`
   - 3.3 Encoding structs
     1. defined a custom `Movie` struct in `internal/data` package.
     2. updated `showMovieHandler` by adding the `Movie` struct
     3. set keys in the `Movie` struct
     4. hid the `CreatedAt` field in the JSON struct using the hyphen (-) which makes the output for this field not to show, set the `Year`, `Runtime`, and `Genre` fields to omitempty which makes the field not to show if it is empty.
   - 3.4 Formatting and Enveloping Responses
     1. changed from `json.Marshal()` to `MarshalIndent()` for project formatting. (Although you take a performance hit, so use only when the data is not resource intensive).
     2. Enveloped our response by showing a JSON key `movie` when rendering the value. Updated the data value in `writeJSON` to be a type of envelope (a map[string]interface{}). Updated `showMovieHandler` to create an instance of the envelope map containing the movie data. Updated the `healthcheckHandler` function to use the envelope map, and removed the `data` variable. 
   - 3.5 Advanced JSON Customization
     1. created `internal/data/runtime.go` to set the runtime field in the Movie struct to be a string i.e. 102 mins
     2. changed the datatype of runtime in the movie struct from `int` to `Runtime`
   - 3.6 Sending Error Messages
     1. created `errors.go` for managing errors.
     2. replaced all http.notFound(), http.Server() errors with the new error handling methods.
     3. Routing errors: errors from `httprouter` router.
     4. Updated `routes.go` by setting a custom error handler for the routes

4. Parsing JSON Requests
   - 4.1: JSON Decoding
     1. updated `createMovieHandler` by adding struct to hold user input
   - 4.2: Managing Bad Requests
     1. triage the Decode error to have a better, cleaner error returned:
        - created `readJSON` to decode the JSON from the request body, triage the errors and replace them with our custom messages.
        - updated `createMovieHandler` by adding the `readJSON` method
     2. updated `createMovieHandler` by adding a newly created method `badRequestResponse` to return 400 error code.
   - 4.3: Restricting Inputs
     1. set limit for the size of the request body in  `readJSON`
     2. set `DisallowUnknownFields` to disallow unknown fields in the request body
   - 4.4: Custom JSON Decoding
     1. updated input struct in `createMovieHandler` to use the `data.Runtime` type for the Runtime field
     2. added `UnmarshalJSON` method to `runtime.go` which will convert the runtime number to int32, and assign to the `Runtime` value itself
   - 4.5: Validating JSON input
     1. defined a custom `Validator` package which contains a map of errors. Methods include: `Check()` which conditionally adds errors to the map, `Valid()` which returns wether the errors map is empty or not. Other functions include `In()`, `Matches()`, and `Unique()` to perform specific validation checks.
     2. added `failedValidationResponse` to return `status unprocessable entity (error 422)` to `errors.go` 
     3. updated `createMovieHandler` by adding the validation pacakage.

5. Database Setup and Configuration
   - 5.1: Setting up PostgreSQL
     1. installed postgres
     2. used `psql` command to create a new database called `greenlight`: `CREATE DATABASE greenlight`
     3. create user `greenlight` without superuser permissions, with a password-based authentication: `CREATE ROLE greenlight WITH LOGIN PASSWORD 'greenlight';`
     5. create a Postgres extension `citext` which adds case-insensitive character string type: `CREATE EXTENSION IF NOT EXISTS citext;`
     6. Exited: `exit` and connected as the `greenlight user`: `psql --host=localhost --dbname=greenlight --username=greenlight`
     7. checked current user: `SELECT current_user;`, check the location of the `postgresql.conf`: `sudo -u postgres psql -c 'SHOW config_file;'`
   - 5.2: Connecting to PostgreSQL
     1. installed postgres [`lib/pq`](github.com/lib/pq) driver
     2. added `openDB` function to connect to the PostgreSQL database
     3. added `.env` for connecting to db using the [Godotenv](github.com/joho/godotenv)
   - 5.3: Configuring the database connection pool
     1.  added `maxOpenConns`, `maxIdleConns`, `maxIdleTime` to the config struct, set the values to 25, 25 and 15 minutes respectively in `openDB` function

6. SQL Migrations
   - 6.1: An overview of SQL Migrations
     1. installed the golang-`migrate` tool using brew i.e. `brew install golang-migrate`
   - 6.2: Working with SQL Migrations
     1. created new migration file for movies table using the `migrate` command: `migrate create -seq -ext=.sql -dir=./migrations create_movies_table`
     2. added the create table movies sql query to the `.up.sql` migration file and drop table movies sql query to the `.down.sql` migration file
     3. created new migration file for containing `CHECK` constraints to ensure business rules at the db level: `migrate create -seq -ext=.sql -dir=./migrations add_movies_check_constraints`. Added the corresponding queries to the `.up.sql` and `.down.sql` files.
     4. added `Makefile` to perform the `migrateup`

7. CRUD Operations
   - 7.1: Setting up the Movie Model
     1. Set up the `model` or `data access/storage` layer which encapsulates all the code for reading and writing the movie data to and from Postgres db.
     2. updated `internal\data\movies.go` by adding the `MovieModel` struct, and CRUD methods for manipulation 
     3. wrapped `MovieModel` in the `Models` struct in a new file `internal/data/models.go`. This is optional but it gives a 'single container' which can hold and represent all the database models as the application grows.
     4. updated `main.go` by adding the models field to `application` struct.
     5. Note, with a few tweak, mocking the database is easy.
   - 7.2: Creating a New Movie
     1. updated `Insert` method in the `internal/data/movies.go` by adding the SQL query and execution statement
     2. updated `createMovieHandler` in `cmd/api/movies.go` by adding the updated `Insert` method
   - 7.3: Fetching a movie
     1. updated the `Get` method by adding the SQL query and execution statement.
     2. updated `showMovieHandler` in `cmd/api/movies.go` by adding the updated `Insert` method
   - 7.4: Updating a movie
     1. added the rotutes for updating the movie
     2. updated the `Update` method by adding the SQL query and execution statement.
     3. created `updateMovieHandler` in `cmd/api/movies.go` and added the updated `Update` method
   - 7.5: Deleting a movie
     1. added the rotutes for deleting the movie
     2. updated the `Delete` method by adding the SQL query and execution statement.
     3. created `deleteMovieHandler` in `cmd/api/movies.go` and added the updated `Update` method

8. Advanced CRUD Operations
   - 8.1: Handling Partial Updates
     1. changed the datatypes of the `input` struct in `updateMovieHandler` method in the `cmd/api/movies.go` file from pass-by-values to pointers (pass-by-reference).
     2. updated the route method from `PUT` to `PATCH` because it is best to use Patch for partial updates on a resource, rather than Put (which is intended for replacing a resource in full)
   - 8.2: Optimistic Concurrency Control
     1. Prevent data race using optimistic locking vs pessimistic locking
     2. created custom `ErrEditConflict` error in case of a conflict
     3. updated `Update` method in `internal/data/movies.go`
     4. added `editConflictResponse` method to `cmd/api/errors.go`
     5. added `editConflictResponse` to `updateMovieHandler`
     6. added the `version` header check in `updateMovieHandler` which allows the sent request to verify that the movie version in the database matches the expected version specified in the header.
   - 8.3: managing SQL Query Timeouts
     1. added context to `Get`, `Update`, `Insert` and `Delete` methods which returns error for time out

9. Filtering, Sorting and Pagination
   - 9.1: Parsing Query String Parameters
     1. created helper functions `readString()`, `readInt()` and `readCSV()` to extract and parse values from query string or return a default 'fallback' value if necessary.
     2. added the new route and updated the code for `listMoviesHandler` to get all movies 
     3. created `filter.go` to handle page filters, moved `PageSize, Page and Sort` fields to `Filters` struct and added `Filters` struct as field to the `input` struct in the `listMoviesHandler` method
   - 9.2: Validating Query String Parameters
     1. added `ValidateFilters` function to certify that the page value is between 1 and 10,000,000; the page_size value is between 1 and 100; the sort parameter contains either "id", "title", "year", "runtime", "-id", "-title", "-year" or "-runtime", where those with `-` infront meaning descending order.
     2. updated `listMoviesHandler` by setting the supported values in the `SortSafelist` field.
   - 9.3: Listing Data
     1. created `GetAll()` method to execute db query
     2. updated `listMoviesHandler` by adding `GetAll()`.
     
      **Note: PostgreSQL features - 1**<br/>
      For dynamic filtering, Postgres has full-text search query features, which can also allow each filter to behave like it is 'optional'. For example:
        - `LOWER(title) = LOWER($1) OR $1 = ''` will evaluate as true if the placeholder parameter $1 is a case-insensitive match for the movie title or the placeholder parameter equals ''. So this filter condition will essentially be ‘skipped’ when movie title being searched for is the empty string "".
         - `genres @> $2 OR $2 = '{}'` condition works in the same way. The @> symbol is the 'contains' operator for PostgreSQL arrays, and this condition will return true if each value in the placeholder parameter $2 appears in the database genres field or the placeholder parameter contains an empty array
   - 9.4: Filtering Lists
     1. updated `GetAll()` method by adding filter query string
   - 9.5: Full-Text Search
     1. updated the query in the `GetAll()` method
     2. created new migration files for movies indexes and added the contents for both up and down: `migrate create -seq -ext .sql -dir ./migrations add_movies_indexes`. Note: We are creating [GIN indexes](https://www.postgresql.org/docs/current/textsearch-indexes.html) on both the _genres_ field and the lexemes generated by `to_tsvector()`, both which are used in the `WHERE` clause of the SQL query. Ran the migration file: `make migrateup`.
      
        **Note: PostgreSQL features - 2**<br/>
        Postgres offers full-text search functionality which allows us to perform 'natural language' searches on the text fields in the database. For example:
         - `to_tsvector('simple', title)` function takes a movie title and splits it into lexemes. By 'simple', the lexemes are lowercase versions of the words in the title, e.g. the movie title = "The Breakfast Club" will be split into lexemes 'breakfast' 'club' 'the'.
         - `plainto_tsquery('simple', $1)` function takes a search value and turns it into a formatted query term that Postgres full-text search can understand. It normailzes the search value (using the 'simple' configuration), strips any special characters and inserts the AND i.e. & operator between the words. Example: search value = "The Club", the query term would be = 'the' & 'club'
         - the `@@` operator is the _matches_ operator. Used to check that the query term matches the lexemes. 
         - A complete list can be found [here](https://www.postgresql.org/docs/9.6/functions-array.html)

   - 9.6: Sorting Lists
     1. client controls the sort order via a query string parameter in the format `sort={-}{field_name}`, where the optional `-` is used to indicate descending order. E.g. `localhost:4000/v1/movies?sort=title`
     2. updated `Filters` struct in `filters.go` by including sortColumn() and sortDirection() helpers to convert query string value (.e. -title) to usable values in our SQL query.
     3. added the newly created method helpers to `internal/data/movies.go`, and updated the SQL query by adding _ORDER BY_ clause.
   - 9.7: Paginating Lists
     1. added `limit()` and `offset()` helper methods to Filters struct for calculating the limit and offset (offset = (page - 1) * page_size) values. **Note:** there is a theoretical risk of integer overflow due to the multiplication of two int values together. However, this is mitigated by the validation rules created in the `ValidateFilters()` function, where the maximum values of page_size=100 and page=10000000 (10 million) were enforced. This means the value returned by `offset()` should never come close to overflowing.
     2. added the created helper methods, and updated the SQL query to the `internal/data/movies.go`
   - 9.8: Returning Pagination Metadata
     1. created `Metadata` struct in `filters.go` to hold the pagination metadata, `calculateMetaData()` function to calculate the appropriate metadata values
     2. added `count(*) OVER()` to SQL query in `internal/data/movies.go` to count the total number of movies (qualifying rows)
     3. updated `listMoviesHandler` to accept metadata struct and values

10. Structured Logging and Error Handling
   - 10.1: Structured JSON Log Entries
     1. Each log entry will be single JSON object containing level, time, message, properties (key/value pairs - optional) and trace (optional) e.g. `{"level":"INFO","time":"2020-12-16T10:53:35Z","message":"starting server","properties":{"addr":":4000","env":"development"}}`
     2. created a custom logger with structs and helper methods in `jsonlog.go` 
     3. updated `main.go` to use the newly created `jsonlog.Logger` struct as a field for the logger, updated the logger variable, and log.Print or log.Fatal commands.
     4. updated `errors.go`, `healthcheck.go` and `movies.go` logger instances.
   - 10.2: Panic Recovery
     1. created `middleware.go` and added `recoverPanic()` method to recover from a panic.
     2. Added the `recoverPanic()` method to `routes.go` which returns http.Handler instead of *httprouter.ROuter

11. Rate Limiting
   - 11.1: Global Rate Limiting
     1. Global rate limiting is useful for enforcing a strict limit on the total rate of requests to the API and you don't care where the requests are coming from.
     2. installed [rate package](golang.org/x/time/rate@latest)
     3. How token-bucket rate limiters work:
        1. A limiter controls how frequently eventsare allowed to happen. It implements a "token bucket" of size _b_, initially full and refilled at rate _r_ tokens per second.
        2. In the context of our API application: 
           - a bucket starting with _b_ tokens.
           - each time a HTTP request is received, we remove one token from the bucket.
           - every 1/r seconds, a token is added back to the bucket - up to a maximum of _b_ total tokens.
           - if we receive a HTTP request and the bucket is empty, return _429 Too Many Requests_ response.
     4. created `rateLimit()` middleware method which creates a new rate limiter for every request that it subsequently handles.
     5. updated `errors.go` by adding the `rateLimitExceededResponse()` methods and added the `rateLimit` middleware to `routes.go` 
   - 11.2: IP-based Rate Limiting
     1. IP-based routing is more common and is used to separate rate limiter for each client, so a client making too many requests does not affect the other clients and their requests.
     2. A conceptual way of implementing is to create an in-memory _map of rate limiters_ using the IP address for each client as the map key.
     3. How it works:
        - Each time a new client makes a request to the API, a new rate limiter is initialized and added to the map.
        - For subsequent requests, the client's rate limiter is retreived from the map and the request is checked if it's permitted by calling its `Allow()` method.
        - Due to the potential of having multiple goroutines accessing the map concurrently, we need to protect access to the map by using `mutex` to prevent race conditions.
     4. updated `rateLimit()` method to used IP based routing. Also add `last seen` feature to prevent the map from growing indefinitely and taking up resources. <br/>
     **Note:** This pattern works for the application running on a sing;e machine, if your infrastructure is distributed, you will need an alternative approach. Alternatively, redis can be used to maintain a request count for clients, running on a server which all the application servers can communicate with. <br/>
   - 11.3: Configuring the Rate Limiters
     1. Make rate limiting values i.e. requests-per-second and burst values) easily configurable so Rate limiting can be turned off if carrying out benchmarks or load testing.
     2. updated `main.go`, and `rateLimit()` method

12. Graceful Shutdown
    1.  _How to safely stop your running application_
    - 12.1: Sending Shutdown Signals
      1. Common signals:
        - Signal -> Description -> Keyboard shortcut -> Catchable
          <br/>SIGINT -> Interrupt from keyboard -> Ctrl+C -> Yes
          <br/>SIGQUIT -> Quit from keyboard -> Ctrl+\ -> Yes
          <br/>SIGKILL -> Kill process (terminate immediately) -> - (none) -> No
          <br/>SIGTERM -> Terminate process in orderly manner -> - (none) -> Yes
        - Catachable signals can be intercepted by our application and either ignored, or used to trigger a certain action (such as a graceful shutdown).
        - How to use: run your server in termina l 1, in terminal 2, get the process ID using `pgrep -l <server_name>` e.g. `pgrep -l api`, finally, terminate the application using `pkill -<signal> <server_name>` e.g. `pkill -SIGINT api`
    - 12.2: Intercepting Shutdown Signals
      1. moved the `http.Server` code to a new file `server.go`
    - 12.3: Executing the shutdown
      1. `Shutdown()` works by first closing all open listeners, then closing all idle connections, and then waiting indefinitely for connections to return to idle and then shut down.
      2. updated `server.go` to receive a SIGINT or SIGTERM signal, which instructs the server to stop accepting any new HTTP requests, and give any in-flight requests a ‘grace period’ of 5 seconds to complete before the application is terminated.
      3. added `time.sleep` to `healthcheckHandler` to test out the functionality.

13. User Model Setup and Registration
    - 13.1: Setting up the Users Database Table
      1. created migration file for creating users: `migrate create -seq -ext=.sql -dir=./migrations create_users_table`
      2. added the SQL queries to the migration files, and ran the migrate command to make migrations: `migrate -path=./migrations -database=$GREENLIGHT_DB_DSN up`
    - 13.2: Setting up the Users Model
      1. created `internal/data/users.go` for holding `User` struct for an individual user, and `UserModel` type for performing SQL queries on the users table.
      2. Added the `bcrypt` package: `go get golang.org/x/crypto/bcrypt@latest`
      3. added validation checks for email, plaintext passwords and users in `internal/data/users.go`
      4. added methods such as Insert, GetByEmail and Update for the user model.
     - 13.3: Registering a User