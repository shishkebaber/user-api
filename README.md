# User API
Backend developer tech test

To run integration tests with postrgesDB execute command in project root folder:
   - docker-compose -f docker-compose.test.yml up --build --abort-on-container-exit

To run application execute command in project root folder:
- docker-compose up --build

By default API runs on :9090 port. Please use http://localhost:9090/ address for local testing. You can find some docs at http://localhost:9090/docs 

By default gRPC server runs on :9092 port

PostgresDB was chosen because of possible future implementation of Users relations, such as friendship or subscription.
Using pgxpool guarantees concurrency-safe data operations. 

We can easily switch to another data base, if needed.

Improvements TODO:
   - Session management (HTTPS, middleware, JWT)
   - DB Migrations improvements
   - Some documentation improvements
   - Pagination 
   - Health Checks
   - Tests improvements (Add GRPC to integration tests)
   - Dockerfile improvements
   - Add better configuration
   
