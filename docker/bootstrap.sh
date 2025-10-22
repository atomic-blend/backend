#!/bin/bash

# remplace restApiUrl inside /usr/share/nginx/html by generated address from env

# Check that PUBLIC_ADDRESS is set
if [ -z "$PUBLIC_ADDRESS" ]; then
    echo "Error: PUBLIC_ADDRESS environment variable is not set"
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