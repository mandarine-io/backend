# Deployment

In this directory we store all files and scripts that need to be deployed to local, development, test or production
environments.

<a id="folder-structure"></a>
## Folder structure

```shell
.
├── README.md
├── config # Configuration files for different tools
│   └── nginx
│       └── nginx.conf
├── docker # Docker compose files for different environments
│   ├── .env.dev
│   ├── .env.test
│   ├── docker-compose.dev.yml
│   ├── docker-compose.local.yml
│   └── docker-compose.test.yml
└── images # Docker image files for building different useful images
    ├── app.Dockerfile
    ├── postgis.Dockerfile
    └── psql.Dockerfile
```

> **Note:** YAML config file and environment file for **application** are located in `<root>/config` directory.

<a id="local-deployment"></a>
## Local deployment

Local environment requires `docker-compose.local.yml` compose file to be deployed. It contains PostgreSQL, Redis and
Minio containers. It use for local development.

Run local environment can be done with:

```bash
make deploy.local
```

or

```bash
cd deploy/docker
cp ../../config/.env.example .env.local
nano .env.local
docker compose -f docker-compose.local.yml --env-file .env.local up -d
```

<a id="development-deployment"></a>
## Development deployment

Development environment requires `docker-compose.dev.yml` compose file to be deployed. It contains PostgreSQL, Redis,
Minio, Nginx and Backend containers. It use for deployment development environment for other developers (WEB and
mobile).

Run development environment can be done with:

```bash
make deploy.dev
```

or

```bash
cd deploy/docker
cp ../../config/.env.example .env.dev
nano .env.dev
docker compose -f docker-compose.dev.yml --env-file .env.dev up -d
```

<a id="test-deployment"></a>
## Test deployment

Test environment requires `docker-compose.test.yml` compose file to be deployed. It contains PostgreSQL, Redis, Minio,
Mailhog, Backend containers and Test data migration job. It use for load testing.

Run test environment can be done with:

```bash
make deploy.test
```

or

```bash
cd deploy/docker
cp ../../config/.env.example .env.test
nano .env.test
docker compose -f docker-compose.test.yml --env-file .env.test up -d
```