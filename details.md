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