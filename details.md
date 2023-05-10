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

3. Chapter 3: Sending JSON Responses
   - 3.0: Sending JSON Responses
     1. updated API handlers to return JSON responses instead of plain text