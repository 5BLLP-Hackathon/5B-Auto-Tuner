#!/bin/bash

# Configuration
VERSION="latest"
LOG_LEVEL="INFO"
DB_NAME="crapi"
DB_USER="admin"
DB_PASSWORD="crapisecretpassword"
DB_HOST="localhost"
DB_PORT="5432"
IDENTITY_SERVER_PORT="8080"
COMMUNITY_SERVER_PORT="8087"
WORKSHOP_SERVER_PORT="8000"
WEB_PORT_HTTP="8888"
WEB_PORT_HTTPS="8443"
MONGO_LOCAL_CONN_URL=mongodb+srv://droy:1HHW1Hp5mGZ9pTSD@5bcluster.kqnuvz6.mongodb.net/?retryWrites=true&w=majority&appName=5BCluster
MONGO_DB_NAME=5bcluster
MAILHOG_HOST="localhost"
MAILHOG_PORT="1025"
MAILHOG_WEB_PORT="8025"
SMTP_HOST="smtp.example.com"
SMTP_PORT="587"
SMTP_EMAIL="user@example.com"
SMTP_PASS="xxxxxxxxxxxxxx"
SMTP_FROM="no-reply@example.com"
SMTP_AUTH="true"
SMTP_STARTTLS="true"
JWT_SECRET="crapi"
JWT_EXPIRATION="604800000"
API_GATEWAY_URL="https://api.mypremiumdealership.com"
TLS_ENABLED="false"
TLS_KEYSTORE_TYPE="PKCS12"
TLS_KEYSTORE="classpath:certs/server.p12"
TLS_KEYSTORE_PASSWORD="passw0rd"
TLS_KEY_PASSWORD="passw0rd"
TLS_KEY_ALIAS="identity"
ENABLE_SHELL_INJECTION="false"
ENABLE_LOG4J="false"

# Function to start PostgreSQL
start_postgresql() {
  echo "Starting PostgreSQL..."
  sudo systemctl start postgresql
  echo "Configuring PostgreSQL..."
  sudo -u postgres psql -c "DO \$$ BEGIN IF NOT EXISTS (SELECT FROM pg_catalog.pg_roles WHERE rolname = '$DB_USER') THEN CREATE ROLE $DB_USER WITH LOGIN PASSWORD '$DB_PASSWORD'; END IF; END\$$;"
  sudo -u postgres psql -c "DO \$$ BEGIN IF NOT EXISTS (SELECT FROM pg_database WHERE datname = '$DB_NAME') THEN CREATE DATABASE $DB_NAME; END IF; END\$$;"
  sudo -u postgres psql -c "GRANT ALL PRIVILEGES ON DATABASE $DB_NAME TO $DB_USER;"
}

# Function to start MailHog
start_mailhog() {
  echo "Starting MailHog..."
  nohup MailHog &>/dev/null &
}

# Function to start crapi-identity
start_crapi_identity() {
  echo "Starting crapi-identity..."
  if [[ -f "crapi-identity-${VERSION}.jar" ]]; then
    java -jar crapi-identity-${VERSION}.jar \
      --LOG_LEVEL=$LOG_LEVEL \
      --DB_NAME=$DB_NAME \
      --DB_USER=$DB_USER \
      --DB_PASSWORD=$DB_PASSWORD \
      --DB_HOST=$DB_HOST \
      --DB_PORT=$DB_PORT \
      --SERVER_PORT=$IDENTITY_SERVER_PORT \
      --ENABLE_SHELL_INJECTION=$ENABLE_SHELL_INJECTION \
      --JWT_SECRET=$JWT_SECRET \
      --MAILHOG_HOST=$MAILHOG_HOST \
      --MAILHOG_PORT=$MAILHOG_PORT \
      --MAILHOG_DOMAIN=$SMTP_EMAIL \
      --SMTP_HOST=$SMTP_HOST \
      --SMTP_PORT=$SMTP_PORT \
      --SMTP_EMAIL=$SMTP_EMAIL \
      --SMTP_PASS=$SMTP_PASS \
      --SMTP_FROM=$SMTP_FROM \
      --SMTP_AUTH=$SMTP_AUTH \
      --SMTP_STARTTLS=$SMTP_STARTTLS \
      --JWT_EXPIRATION=$JWT_EXPIRATION \
      --ENABLE_LOG4J=$ENABLE_LOG4J \
      --API_GATEWAY_URL=$API_GATEWAY_URL \
      --TLS_ENABLED=$TLS_ENABLED \
      --TLS_KEYSTORE_TYPE=$TLS_KEYSTORE_TYPE \
      --TLS_KEYSTORE=$TLS_KEYSTORE \
      --TLS_KEYSTORE_PASSWORD=$TLS_KEYSTORE_PASSWORD \
      --TLS_KEY_PASSWORD=$TLS_KEY_PASSWORD \
      --TLS_KEY_ALIAS=$TLS_KEY_ALIAS &
  else
    echo "Error: crapi-identity-${VERSION}.jar not found."
  fi
}

# Function to start crapi-community
start_crapi_community() {
  echo "Starting crapi-community..."
  if [[ -f "crapi-community-${VERSION}.jar" ]]; then
    java -jar crapi-community-${VERSION}.jar \
      --LOG_LEVEL=$LOG_LEVEL \
      --IDENTITY_SERVICE=http://localhost:$IDENTITY_SERVER_PORT \
      --DB_NAME=$DB_NAME \
      --DB_USER=$DB_USER \
      --DB_PASSWORD=$DB_PASSWORD \
      --DB_HOST=$DB_HOST \
      --DB_PORT=$DB_PORT \
      --SERVER_PORT=$COMMUNITY_SERVER_PORT \
      --MONGO_DB_HOST=$MONGO_LOCAL_CONN_URL \
      --MONGO_DB_NAME=$MONGO_DB_NAME \
      --TLS_ENABLED=$TLS_ENABLED \
      --TLS_CERTIFICATE=certs/server.crt \
      --TLS_KEY=certs/server.key &
  else
    echo "Error: crapi-community-${VERSION}.jar not found."
  fi
}

# Function to start crapi-workshop
start_crapi_workshop() {
  echo "Starting crapi-workshop..."
  if [[ -f "crapi-workshop-${VERSION}.jar" ]]; then
    java -jar crapi-workshop-${VERSION}.jar \
      --LOG_LEVEL=$LOG_LEVEL \
      --IDENTITY_SERVICE=http://localhost:$IDENTITY_SERVER_PORT \
      --DB_NAME=$DB_NAME \
      --DB_USER=$DB_USER \
      --DB_PASSWORD=$DB_PASSWORD \
      --DB_HOST=$DB_HOST \
      --DB_PORT=$DB_PORT \
      --SERVER_PORT=$WORKSHOP_SERVER_PORT \
      --MONGO_DB_HOST=$MONGO_LOCAL_CONN_URL \
      --MONGO_DB_NAME=$MONGO_DB_NAME \
      --SECRET_KEY=$JWT_SECRET \
      --API_GATEWAY_URL=$API_GATEWAY_URL \
      --TLS_ENABLED=$TLS_ENABLED \
      --TLS_CERTIFICATE=certs/server.crt \
      --TLS_KEY=certs/server.key &
  else
    echo "Error: crapi-workshop-${VERSION}.jar not found."
  fi
}

# Function to start crapi-web
start_crapi_web() {
  echo "Starting crapi-web..."
  nohup npm start --prefix ./js &>/dev/null &
}

# Start services
start_postgresql
start_mailhog
start_crapi_identity
start_crapi_community
start_crapi_workshop
start_crapi_web

echo "All services are up and running."

# Wait for services to be stopped
wait
