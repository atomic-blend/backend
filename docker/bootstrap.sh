#!/bin/bash

# remplace restApiUrl inside /usr/share/nginx/html by generated address from env

# Check that PUBLIC_ADDRESS is set
if [ -z "$PUBLIC_ADDRESS" ]; then
    echo "Error: PUBLIC_ADDRESS environment variable is not set"
    exit 1
fi

# The first positional argument is the app name (e.g. mail, task)
APP_NAME="$1"
if [ -z "$APP_NAME" ]; then
    echo "Usage: $0 <app-name>"
    echo "Example: $0 mail"
    exit 1
fi

# Generate the rest API URL based on HTTPS setting
# Check if PUBLIC_ADDRESS is an IP address (IPv4 pattern)
if [[ "$PUBLIC_ADDRESS" =~ ^[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
    # It's an IP address, don't add api. prefix
    if [ "$HTTPS" = "true" ]; then
        REST_API_URL="https://${PUBLIC_ADDRESS}"
    else
        REST_API_URL="http://${PUBLIC_ADDRESS}"
    fi
else
    # It's a domain name, add api. prefix
    if [ "$HTTPS" = "true" ]; then
        REST_API_URL="https://api.${PUBLIC_ADDRESS}"
    else
        REST_API_URL="http://api.${PUBLIC_ADDRESS}"
    fi
fi

echo "Generated REST API URL: $REST_API_URL"

# Replace the restApiUrl in the assets/assets/configs/prod.json
sed -i "s|\"restApiUrl\":[[:space:]]*\"[^\"]*\"|\"restApiUrl\": \"$REST_API_URL\"|g" /usr/share/nginx/html/assets/assets/configs/prod.json

echo "Updated prod.json with REST API URL"

# Build public URL for the given app name and update publicUrl key in prod.json
# If PUBLIC_ADDRESS is an IP, we won't prefix with the app name (same behaviour as restApiUrl logic)
if [[ "$PUBLIC_ADDRESS" =~ ^[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
    if [ "$HTTPS" = "true" ]; then
        PUBLIC_URL="https://${PUBLIC_ADDRESS}"
    else
        PUBLIC_URL="http://${PUBLIC_ADDRESS}"
    fi
else
    if [ "$HTTPS" = "true" ]; then
        PUBLIC_URL="https://${APP_NAME}.${PUBLIC_ADDRESS}"
    else
        PUBLIC_URL="http://${APP_NAME}.${PUBLIC_ADDRESS}"
    fi
fi

echo "Generated PUBLIC URL for app '$APP_NAME': $PUBLIC_URL"

# Replace the publicUrl in the assets/assets/configs/prod.json
sed -i "s|\"publicUrl\":[[:space:]]*\"[^\"]*\"|\"publicUrl\": \"$PUBLIC_URL\"|g" /usr/share/nginx/html/assets/assets/configs/prod.json

echo "Updated prod.json with PUBLIC URL for app '$APP_NAME'"