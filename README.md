# user-service

## Description

REST endpoints for user management:
- POST /users - Create a new user.
- GET /users/{id} - Retrieve user details by ID.
- PUT /users/{id} - Update user information by ID.
- DELETE /users/{id} - Delete a user by ID.
- GET /users - List all users with pagination.

Repository: https://github.com/bmcszk/user-service

## Setup instructions

### Prerequisites
- Docker
- Docker Compose
- Git
- Goland v1.23 or higher
- Make
- (optional) Kubectl

### Installation options

- `make compose` - starts the docker compose environment
- `make test` - runs tests using compose
- `make coverage` - creates coverage report
- `make kind-up tilt-up` - starts local [Kind](https://kind.sigs.k8s.io/) k8s cluster and [Tilt](https://docs.tilt.dev/) deployment

## Development
- `make sqlc` - generates code using sqlc
- `make migrate-up` - runs migrations up
- `make migrate-down` - runs migrations down

## Design decisions

1. Simple project structure, directories:
- sources:
    - `db` - database layer
    - `api` - api layer
    - `logic` - business logic layer
    - `e2e` - end-to-end tests
- configs:
    - `helm` - helm charts
    - `tilt` - tilt configs
2. Simplest http layer using standard library
3. Database is [Postgres](https://www.postgresql.org/) 
4. Simplest db layer using standard library and [SQLC](https://sqlc.dev/) for code generation.
5. Migrations handled automatically using [Golang Migrate](https://github.com/golang-migrate/migrate)
6. Deployment can be managed by Helm
7. Local dev environment can be managed by [Tilt](https://docs.tilt.dev/)
8. CI is done using Github Actions where Kind cluster is created and Tilt is deployed. Check: https://github.com/bmcszk/user-service/actions/runs/11915468146/job/33205800016

## Assumptions
1. User model is very basic but with possibility to add more fields through DB migrations
2. No authentication, no authorization
3. Many TODOs in code

## Scalability Considerations

1. Application can easily scaled using K8s manual and automatical scaling. Database can easily scaled using [Cloud Native Postgres](https://cloudnative-pg.io/)
2. Designing the system to be able to handle 30k requests per minute with a latency of < 100ms. Step by step:
    1. Prepare configurable e2e test scenario running with up to 30k requests per minute.
    2. Prepare configuration and middlewares for telemetry, to monitor performance and latency of application and DB.
    3. Prepare test K8s cluster and deploy the application with DB.
    4. Run the e2e test scenario with small amount of requests per minute.
    5. Monitor the performance and latency of the application and DB.
    6. Gradually increase the number of requests per minute and monitor the performance and latency of the application and DB.
    7. Identify bottlenecks and adjust resources configuration, adjust number of instances of the application and DB.
    8. Go to step 5. until reaching the desired performance and latency.

## Security

1. Basic input validation introduced.
2. DB queries are sanitized using prepared statements with params.
3. Introducing SSL strongly recommended.
3. Authentication and authorization can be implemented in Kubernetes using 3rd party solution like [Kong](https://konghq.com/) or [Nginx](https://www.nginx.com/).
