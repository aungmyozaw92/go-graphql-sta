# Installation

```bash
$ go mod tidy
```

## Running the app

```bash
# development
$ go run main.go or go run .

```

## Database Configuration

```bash

$ DB_HOST=127.0.0.1
$ DB_USER=
$ DB_PASSWORD=
$ DB_NAME=
$ DB_PORT=

```

## Redis Configuration

```bash

$ REDIS_ADDRESS=127.0.0.1:6379
$ CACHE_LIFESPAN=4

```

## API Configuration

```bash

$ API_PORT=8080
$ API_SECRET=your_secret_key
$ TOKEN_HOUR_LIFESPAN=240

```

## Storage Configuration

```bash
#Storage Service Configuration

#Options: DO_SPACE for Digital Ocean Space, GOOGLE_CLOUD for Google Cloud Service

$ STORAGE_SERVICE=GOOGLE_CLOUD or
$ STORAGE_SERVICE=DO_SPACE

# Digital Ocean Space Config

$ SP_ACCESS_KEY_ID=
$ SP_SECRET_ACCESS_KEY=
$ SP_DEFAULT_REGION=
$ SP_BUCKET=
$ SP_URL=


# Google Cloud Service Config

$ GOOGLE_CLOUD_CREDENTIALS='{
"type": "service_account",
"project_id": "your_project_id",
...
}'
$ GCS_BUCKET=your_gcs_bucket
$ GCS_URL=storage.googleapis.com

```

## GORM Logging

```bash
$ GORM_LOG=gorm.log
```

#
