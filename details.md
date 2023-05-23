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