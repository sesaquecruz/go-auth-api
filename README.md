# An Auth API with JWT

This project is a RESTful API written in Go for user authentication and management. It uses JSON Web Tokens (JWT) for authorization and supports basic CRUD operations for user accounts. The API is designed to be easy to integrate into other projects and provides a clean and efficient way to manage user access.

## Endpoints

| Endpoint | Method | Protected | Description |
| -------- | ------ | --------- | ----------- |
| `/api/v1/login` | POST   | NO  | Authenticate user and receive JWT token |
| `/api/v1/users` | POST   | NO  | Create a new user account               |
| `/api/v1/users` | GET    | YES | Retrieve user data                      |
| `/api/v1/users` | PUT    | YES | Update user data                        |
| `/api/v1/users` | DELETE | YES | Delete user account                     |
| `/api/v1/docs`  | GET    | NO  | API Documentation / Swagger UI                              |

## Requirements

To use this program, you will need:

- Docker
- A stable internet connection

## Installation

1. Clone this repository:

```
git clone https://github.com/sesaquecruz/go-auth-api
```

2. Enter the project directory:

```
cd go-auth-api
```

3. Run the docker compose:

```
docker compose up --build
```

## Usage

### API Documentation

1. Access the Swagger UI:

```
http://localhost:8080/api/v1/docs
```

To access protected endpoints, a JWT token is required. This token can be obtained by creating a user account and authenticating it. The user ID will be included in the token.

## Troubleshooting

See [docker-compose.yml](./docker-compose.yml) to verify or change the services, port values, or environment variables values.

## Contributing

Contributions are welcome! If you find a bug or would like to suggest an enhancement, please make a fork, create a new branch with the bugfix or feature, and submit a pull request.

This project follows the [GitFlow](https://www.atlassian.com/git/tutorials/comparing-workflows/gitflow-workflow) workflow and adheres to [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/). It also has a simple [GitHub Action](./.github/workflows/build-and-test.yml) to verify the build and run tests before merging in some branches.

## License

This project is licensed under the MIT License. See [LICENSE](./LICENSE) file for more information.