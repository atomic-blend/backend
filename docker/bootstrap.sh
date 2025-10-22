#!/bin/bash

# remplace restApiUrl inside /usr/share/nginx/html by generated address from env

# Check that PUBLIC_ADDRESS is set
if [ -z "$PUBLIC_ADDRESS" ]; then
    echo "Error: PUBLIC_ADDRESS environment variable is not set"
    exit 1
fi

# Generate the rest API URL based on HTTPS setting
if [ "$HTTPS" = "true" ]; then
    REST_API_URL="https://api.${PUBLIC_ADDRESS}"
else
    REST_API_URL="http://api.${PUBLIC_ADDRESS}"
fi

echo "Generated REST API URL: $REST_API_URL"

# Replace the restApiUrl in the assets/assets/configs/prod.json
sed -i "s|\"restApiUrl\":[[:space:]]*\"[^\"]*\"|\"restApiUrl\": \"$REST_API_URL\"|g" /usr/share/nginx/html/assets/assets/configs/prod.json

echo "Updated prod.json with REST API URL"