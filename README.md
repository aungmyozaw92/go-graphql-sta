Please add .env

DB_HOST=127.0.0.1

DB_USER=

DB_PASSWORD=

DB_NAME=

DB_PORT=

REDIS_ADDRESS=127.0.0.1:6379

CACHE_LIFESPAN=4

API_PORT_2=8080

API_SECRET=your secret key

TOKEN_HOUR_LIFESPAN=240

STORAGE_SERVICE=DO_SPACE # for digital ocean space

if you use Digital Ocean's space ,

please add SP_ACCESS_KEY_ID,SP_SECRET_ACCESS_KEY, SP_DEFAULT_REGION, SP_BUCKET and SP_URL.

SP_ACCESS_KEY_ID=

SP_SECRET_ACCESS_KEY=

SP_DEFAULT_REGION=

SP_BUCKET=

SP_URL=

or

STORAGE_SERVICE=GOOGLE_CLOUD # for goocle cloud service

if you use Google Cloud ,

please add GOOGLE_CLOUD_CREDENTIALS,GCS_BUCKET and GCS_URL.

GOOGLE_CLOUD_CREDENTIALS='{

add your google cloud service credentials

}'

GCS_BUCKET=add your google cloud bucket

GCS_URL=storage.googleapis.com

GORM_LOG=gorm.log
