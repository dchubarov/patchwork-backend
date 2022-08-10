#!/usr/bin/env bash

set -e
source /docker-entrypoint-initdb.d/initdb.env

mongo <<EOF
use admin

db.createUser({
    user: '${MONGO_SERVICE_USERNAME}',
    pwd: '${MONGO_SERVICE_PASSWORD}',
    roles: [
        {
            role: "readWrite",
            db: '${MONGO_DATABASE}'
        }
    ]
});

use ${MONGO_DATABASE}
EOF
