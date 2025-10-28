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

echo "Generated REST API URL: $REST_API_URL"

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

# Path to prod.json
PROD_JSON="/usr/share/nginx/html/assets/assets/configs/prod.json"

if [ ! -f "$PROD_JSON" ]; then
    echo "Error: $PROD_JSON not found"
    exit 1
fi

# Prefer jq for safe JSON updates. If not present, fall back to a portable sed-based approach.
if command -v jq >/dev/null 2>&1; then
    # Use jq to set or create both keys atomically
    tmpfile=$(mktemp)
    jq --arg r "$REST_API_URL" --arg p "$PUBLIC_URL" '.restApiUrl = $r | .publicUrl = $p' "$PROD_JSON" > "$tmpfile" && mv "$tmpfile" "$PROD_JSON"
    if [ $? -ne 0 ]; then
        echo "Error: failed to update $PROD_JSON with jq"
        exit 1
    fi
    echo "Updated prod.json with REST API URL and PUBLIC URL (via jq)"
else
    # Fallback: do safe replace-or-insert using sed (portable) and a temp file.
    # Helper to replace a key if present, else insert before closing brace.
    replace_or_insert() {
        key="$1"
        value="$2"
        tmp=$(mktemp)

        if grep -q "\"${key}\"" "$PROD_JSON"; then
            # key exists -> replace its value
            sed "s|\"${key}\"[[:space:]]*:[[:space:]]*\"[^\"]*\"|\"${key}\": \"${value}\"|g" "$PROD_JSON" > "$tmp" && mv "$tmp" "$PROD_JSON"
        else
            # key missing -> decide whether to prefix with a comma depending on existing keys
            if grep -q '"[[:alnum:]_]\+"' "$PROD_JSON"; then
                # file contains other keys -> add a leading comma
                sed -e ':a' -e 'N' -e '$!ba' -e "s/}\s*$/,\n  \"${key}\": \"${value}\"\n}/" "$PROD_JSON" > "$tmp" && mv "$tmp" "$PROD_JSON"
            else
                # empty object -> don't add leading comma
                sed -e ':a' -e 'N' -e '$!ba' -e "s/}\s*$/  \"${key}\": \"${value}\"\n}/" "$PROD_JSON" > "$tmp" && mv "$tmp" "$PROD_JSON"
            fi
        fi
    }

    replace_or_insert "restApiUrl" "$REST_API_URL"
    replace_or_insert "publicUrl" "$PUBLIC_URL"

    echo "Updated prod.json with REST API URL and PUBLIC URL (sed fallback)"
fi