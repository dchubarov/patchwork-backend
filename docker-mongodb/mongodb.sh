#!/usr/bin/env bash

set -e
source .env

DATA_DIR=$(realpath ../build/docker-mongodb)
echo "Using data directory $DATA_DIR"
mkdir -p "$DATA_DIR"

if [[ ! -f "$DATA_DIR/root.password" ]]; then
  rootPwd=$(pwgen -cnsB 16 1)
  echo -n "$rootPwd" >"$DATA_DIR/root.password"
else
  rootPwd=$(cat "$DATA_DIR/root.password")
fi

if [[ ! -f "$DATA_DIR/user.password" ]]; then
  userPwd=$(pwgen -cnsB 9 1)
  echo -n "$userPwd" >"$DATA_DIR/user.password"
fi

function await_db() {
  success=1
  echo "Awaiting database is up and running $1"
  for _ in {1..10}; do
    if docker exec "$1" mongosh -u "$MONGO_ROOT_USERNAME" -p "$rootPwd" --eval "db.runCommand('ping').ok" >/dev/null; then
      success=0
      break
    else
      sleep 2
    fi
  done

  if [ $success -ne 0 ]; then
    echo "Database failed to start"
  fi

  return $success
}

if [[ ! -f "$DATA_DIR/db/.initialized" ]]; then
  echo "Starting init container..."

  tee "$DATA_DIR/initdb.js" >/dev/null <<EOT
db.disableFreeMonitoring()
db.getSiblingDB('admin').createUser({
    user: '${MONGO_SERVICE_USERNAME}',
    pwd: '$userPwd',
    roles: [
        {
            role: 'readWrite',
            db: '${MONGO_DATABASE}'
        }
    ]
})
EOT

  docker run --detach --rm \
    --cidfile "$DATA_DIR/mongo-init.cid" \
    --name mongo-init \
    -e "MONGO_INITDB_ROOT_USERNAME=$MONGO_ROOT_USERNAME" \
    -e "MONGO_INITDB_ROOT_PASSWORD=$rootPwd" \
    -e "MONGO_INITDB_DATABASE=$MONGO_DATABASE" \
    -v "$DATA_DIR/db:/data/db" \
    -v "$DATA_DIR/initdb.js:/docker-entrypoint-initdb.d/initdb.js:ro" \
    "mongo:$MONGO_VERSION"

  if [[ ! -f "$DATA_DIR/mongo-init.cid" ]]; then
    echo "Could not start init container"
    exit 1
  fi

  cid=$(cat "$DATA_DIR/mongo-init.cid")
  if await_db "$cid"; then
    echo "Database initialization completed"
    touch "$DATA_DIR/db/.initialized"
  fi

  echo "Stopping init container $cid"
  docker stop "$cid" >/dev/null
  rm -f "$DATA_DIR/mongo-init.cid"
fi

if [[ ! -f "$DATA_DIR/replset.key" ]]; then
  echo "Generating replica set key file..."
  openssl rand -base64 256 -out "$DATA_DIR/replset.key"
  chmod 600 "$DATA_DIR/replset.key"
fi

echo "Starting container in replication mode..."
docker run --detach \
  --cidfile "$DATA_DIR/mongo.cid" \
  --name mongo \
  -h mongo \
  -v "$DATA_DIR/db:/data/db" \
  -v "$DATA_DIR/replset.key:/etc/mongod/replset.key" \
  -p "${MONGO_PORT:-27017}:${MONGO_PORT:-27017}" \
  "mongo:$MONGO_VERSION" \
  --replSet "$MONGO_REPLICA_SET" \
  --keyFile /etc/mongod/replset.key \
  --bind_ip_all \
  --auth

if [[ ! -f "$DATA_DIR/mongo.cid" ]]; then
  echo "Could not start mongodb container"
  exit 1
fi

cid=$(cat "$DATA_DIR/mongo.cid")
if await_db "$cid"; then
  if ! docker exec "$cid" mongosh -u "$MONGO_ROOT_USERNAME" -p "$rootPwd" --eval "rs.status().ok" >/dev/null; then
    echo "Replica set has not been initialized yet, initiating now..."
    if ! docker exec "$cid" mongosh -u "$MONGO_ROOT_USERNAME" -p "$rootPwd" --eval "rs.initiate()"; then
      echo "Replica set initialization failed"
      exit 1
    fi
  fi

  echo "Database started in container $cid"
fi
