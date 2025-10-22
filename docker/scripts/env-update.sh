#!/bin/bash

# Environment file update functionality for setup script
# This module handles updating .env files with new versions

# Source utils for colors and common functions
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$SCRIPT_DIR/utils.sh"

# Function to check for outdated services and generate update summary
check_outdated_services() {
    local outdated_count=0
    local outdated_services=()

    for i in "${!IMAGES[@]}"; do
        local base_image="${IMAGES[$i]}"
        local env_var="${ENV_VARS[$i]}"
        local current_version=$(get_current_version "$env_var")
        local latest_version=$(get_latest_version "$base_image")
        
        # Count outdated services
        if [ "$latest_version" != "$current_version" ]; then
            outdated_count=$((outdated_count + 1))
            outdated_services+=("${env_var}=$latest_version")
        fi
    done

    # Export for use in other functions
    export OUTDATED_COUNT=$outdated_count
    export OUTDATED_SERVICES=("${outdated_services[@]}")
}

# Function to display update summary
display_update_summary() {
    echo ""
    echo -e "${BLUE}Summary:${NC}"
    echo "=========="

    # Only show update section if there are outdated services
    if [ $OUTDATED_COUNT -gt 0 ]; then
        echo -e "${YELLOW}Add these to your .env.versions file:${NC}"
        echo ""
        
        for service in "${OUTDATED_SERVICES[@]}"; do
            echo "$service"
        done
        
        echo ""
        echo -e "${BLUE}Would you like to automatically update the .env file with these new versions?${NC}"
        read -p "Enter 'y' or 'yes' to proceed: " -r
        echo ""
        
        if [[ $REPLY =~ ^[Yy]([Ee][Ss])?$ ]]; then
            update_env_file
        else
            echo -e "${YELLOW}Update cancelled. You can manually update the .env file with the values shown above.${NC}"
            echo ""
            echo -e "${BLUE}Would you like to configure atomic-blend settings anyway?${NC}"
            read -p "Enter 'y' or 'yes' to proceed: " -r
            echo ""
            
            if [[ $REPLY =~ ^[Yy]([Ee][Ss])?$ ]]; then
                configure_atomic_blend
            fi
        fi
    else
        echo -e "${GREEN}All services are up to date!${NC}"
        echo ""
        echo -e "${BLUE}Would you like to configure atomic-blend settings?${NC}"
        read -p "Enter 'y' or 'yes' to proceed: " -r
        echo ""
        
        if [[ $REPLY =~ ^[Yy]([Ee][Ss])?$ ]]; then
            configure_atomic_blend
        fi
    fi
}

# Function to update the .env file
update_env_file() {
    local env_file="${ENV_FILE:-.env}"
    
    # Get the absolute path to the docker directory
    local docker_dir="$(dirname "$SCRIPT_DIR")"
    
    # Convert relative path to absolute path
    if [[ "$env_file" != /* ]]; then
        env_file="$docker_dir/$env_file"
    fi
    
    echo -e "${YELLOW}Updating .env file...${NC}"
    
    # Create backup of .env file
    cp "$env_file" "${env_file}.backup.$(date +%Y%m%d_%H%M%S)"
    echo -e "${BLUE}Backup created: ${env_file}.backup.$(date +%Y%m%d_%H%M%S)${NC}"
    
    # Update each outdated service in .env file
    for service in "${OUTDATED_SERVICES[@]}"; do
        local env_var=$(echo "$service" | cut -d'=' -f1)
        local new_version=$(echo "$service" | cut -d'=' -f2)
        
        # Check if the variable exists in .env file
        if grep -q "^${env_var}=" "$env_file"; then
            # Get the original line to check for quotes
            local original_line=$(grep "^${env_var}=" "$env_file")
            
            # Check if the original value was quoted (double quotes)
            if [[ $original_line =~ ^${env_var}=\".*\"$ ]]; then
                # Replace with double quotes
                sed -i.tmp "s/^${env_var}=.*/${env_var}=\"${new_version}\"/" "$env_file"
            # Check if the original value was quoted (single quotes)
            elif [[ $original_line =~ ^${env_var}=\'.*\'$ ]]; then
                # Replace with single quotes
                sed -i.tmp "s/^${env_var}=.*/${env_var}='${new_version}'/" "$env_file"
            else
                # Replace without quotes
                sed -i.tmp "s/^${env_var}=.*/${env_var}=${new_version}/" "$env_file"
            fi
            echo -e "${GREEN}Updated ${env_var} to ${new_version}${NC}"
        else
            # Add new variable with quotes (default format)
            echo "${env_var}=\"${new_version}\"" >> "$env_file"
            echo -e "${GREEN}Added ${env_var}=\"${new_version}\"${NC}"
        fi
    done
    
    # Remove temporary file created by sed
    rm -f "${env_file}.tmp"
    
    echo -e "${GREEN}Successfully updated .env file!${NC}"
    
    # Ask for atomic-blend configuration
    configure_atomic_blend
}

# Function to configure atomic-blend environment variables
configure_atomic_blend() {
    local env_file="${ENV_FILE:-.env}"
    
    # Get the absolute path to the docker directory
    local docker_dir="$(dirname "$SCRIPT_DIR")"
    
    # Convert relative path to absolute path
    if [[ "$env_file" != /* ]]; then
        env_file="$docker_dir/$env_file"
    fi
    
    echo ""
    echo -e "${BLUE}Environment Configuration${NC}"
    echo "========================="
    echo ""
    
    # Check and generate SSO_SECRET if needed
    local current_sso_secret=$(grep "^SSO_SECRET=" "$env_file" 2>/dev/null | cut -d'=' -f2 | tr -d '"' | tr -d "'")
    
    if [ -z "$current_sso_secret" ]; then
        echo -e "${YELLOW}SSO_SECRET is empty or missing. Generating a new one...${NC}"
        
        # Generate SSO_SECRET using openssl
        # Use tr -d '\n' for cross-platform compatibility (macOS base64 doesn't have -w0)
        local sso_secret=$(openssl rand 256 | base64 | tr -d '\n')
        
        # Update or add SSO_SECRET in .env file
        if grep -q "^SSO_SECRET=" "$env_file"; then
            # Get the original line to check for quotes
            local original_line=$(grep "^SSO_SECRET=" "$env_file")
            
            # Check if the original value was quoted (double quotes)
            if [[ $original_line =~ ^SSO_SECRET=\".*\"$ ]]; then
                # Replace with double quotes
                sed -i.tmp "s|^SSO_SECRET=.*|SSO_SECRET=\"${sso_secret}\"|" "$env_file"
            # Check if the original value was quoted (single quotes)
            elif [[ $original_line =~ ^SSO_SECRET=\'.*\'$ ]]; then
                # Replace with single quotes
                sed -i.tmp "s|^SSO_SECRET=.*|SSO_SECRET='${sso_secret}'|" "$env_file"
            else
                # Replace without quotes
                sed -i.tmp "s|^SSO_SECRET=.*|SSO_SECRET=${sso_secret}|" "$env_file"
            fi
            echo -e "${GREEN}Updated SSO_SECRET with newly generated value${NC}"
        else
            # Add new variable with quotes (default format)
            echo "SSO_SECRET=\"${sso_secret}\"" >> "$env_file"
            echo -e "${GREEN}Added SSO_SECRET with newly generated value${NC}"
        fi
        
        # Remove temporary file created by sed
        rm -f "${env_file}.tmp"
    else
        echo -e "${GREEN}SSO_SECRET already exists and is not empty${NC}"
    fi
    
    echo ""
    echo -e "${BLUE}Public Domain Configuration${NC}"
    echo "============================"
    echo -e "${YELLOW}PUBLIC_ADDRESS and ACCOUNT_DOMAINS define where your service will be accessible.${NC}"
    echo ""
    
    # Get current value for PUBLIC_ADDRESS
    local current_public_address=$(grep "^PUBLIC_ADDRESS=" "$env_file" 2>/dev/null | cut -d'=' -f2 | tr -d '"' | tr -d "'")
    
    if [ -n "$current_public_address" ]; then
        echo -e "${BLUE}Current value: ${current_public_address}${NC}"
    else
        echo -e "${BLUE}No current value set${NC}"
    fi
    
    echo -e "${BLUE}Example: app.example.com${NC}"
    echo ""
    read -p "Would you like to update the public domain? (y/n): " -r update_domain
    
    if [[ $update_domain =~ ^[Yy]$ ]]; then
        read -p "Enter public domain: " -r public_domain
        
        if [ -n "$public_domain" ]; then
        # Update or add PUBLIC_ADDRESS in .env file
        if grep -q "^PUBLIC_ADDRESS=" "$env_file"; then
            # Get the original line to check for quotes
            local original_line=$(grep "^PUBLIC_ADDRESS=" "$env_file")
            
            # Check if the original value was quoted (double quotes)
            if [[ $original_line =~ ^PUBLIC_ADDRESS=\".*\"$ ]]; then
                # Replace with double quotes
                sed -i.tmp "s|^PUBLIC_ADDRESS=.*|PUBLIC_ADDRESS=\"${public_domain}\"|" "$env_file"
            # Check if the original value was quoted (single quotes)
            elif [[ $original_line =~ ^PUBLIC_ADDRESS=\'.*\'$ ]]; then
                # Replace with single quotes
                sed -i.tmp "s|^PUBLIC_ADDRESS=.*|PUBLIC_ADDRESS='${public_domain}'|" "$env_file"
            else
                # Replace without quotes
                sed -i.tmp "s|^PUBLIC_ADDRESS=.*|PUBLIC_ADDRESS=${public_domain}|" "$env_file"
            fi
            echo -e "${GREEN}Updated PUBLIC_ADDRESS to ${public_domain}${NC}"
        else
            # Add new variable with quotes (default format)
            echo "PUBLIC_ADDRESS=\"${public_domain}\"" >> "$env_file"
            echo -e "${GREEN}Added PUBLIC_ADDRESS=\"${public_domain}\"${NC}"
        fi
        
        # Update or add ACCOUNT_DOMAINS in .env file
        if grep -q "^ACCOUNT_DOMAINS=" "$env_file"; then
            # Get the original line to check for quotes
            local original_line=$(grep "^ACCOUNT_DOMAINS=" "$env_file")
            
            # Check if the original value was quoted (double quotes)
            if [[ $original_line =~ ^ACCOUNT_DOMAINS=\".*\"$ ]]; then
                # Replace with double quotes
                sed -i.tmp "s|^ACCOUNT_DOMAINS=.*|ACCOUNT_DOMAINS=\"${public_domain}\"|" "$env_file"
            # Check if the original value was quoted (single quotes)
            elif [[ $original_line =~ ^ACCOUNT_DOMAINS=\'.*\'$ ]]; then
                # Replace with single quotes
                sed -i.tmp "s|^ACCOUNT_DOMAINS=.*|ACCOUNT_DOMAINS='${public_domain}'|" "$env_file"
            else
                # Replace without quotes
                sed -i.tmp "s|^ACCOUNT_DOMAINS=.*|ACCOUNT_DOMAINS=${public_domain}|" "$env_file"
            fi
            echo -e "${GREEN}Updated ACCOUNT_DOMAINS to ${public_domain}${NC}"
        else
            # Add new variable with quotes (default format)
            echo "ACCOUNT_DOMAINS=\"${public_domain}\"" >> "$env_file"
            echo -e "${GREEN}Added ACCOUNT_DOMAINS=\"${public_domain}\"${NC}"
        fi
        
            # Remove temporary file created by sed
            rm -f "${env_file}.tmp"
        else
            echo -e "${RED}Error: Public domain cannot be empty.${NC}"
        fi
    else
        echo -e "${YELLOW}Skipped PUBLIC_ADDRESS and ACCOUNT_DOMAINS configuration.${NC}"
    fi
    
    # Ask about AUTH_MAX_NB_USER customization
    echo ""
    echo -e "${BLUE}User Account Limit Configuration${NC}"
    echo "=================================="
    echo -e "${YELLOW}AUTH_MAX_NB_USER defines the maximum number of user accounts allowed.${NC}"
    echo ""
    
    # Get current value
    local current_max_users=$(grep "^AUTH_MAX_NB_USER=" "$env_file" 2>/dev/null | cut -d'=' -f2 | tr -d '"' | tr -d "'")
    
    if [ -n "$current_max_users" ]; then
        echo -e "${BLUE}Current value: ${current_max_users}${NC}"
    else
        echo -e "${BLUE}No current value set${NC}"
    fi
    
    echo ""
    read -p "Would you like to update the maximum number of users? (y/n): " -r update_max_users
    
    if [[ $update_max_users =~ ^[Yy]$ ]]; then
        read -p "Enter maximum number of users: " -r max_users
        
        if [ -n "$max_users" ]; then
            # Validate that input is a number
            if [[ "$max_users" =~ ^[0-9]+$ ]]; then
            # Update or add AUTH_MAX_NB_USER in .env file
            if grep -q "^AUTH_MAX_NB_USER=" "$env_file"; then
                # Get the original line to check for quotes
                local original_line=$(grep "^AUTH_MAX_NB_USER=" "$env_file")
                
                # Check if the original value was quoted (double quotes)
                if [[ $original_line =~ ^AUTH_MAX_NB_USER=\".*\"$ ]]; then
                    # Replace with double quotes
                    sed -i.tmp "s|^AUTH_MAX_NB_USER=.*|AUTH_MAX_NB_USER=\"${max_users}\"|" "$env_file"
                # Check if the original value was quoted (single quotes)
                elif [[ $original_line =~ ^AUTH_MAX_NB_USER=\'.*\'$ ]]; then
                    # Replace with single quotes
                    sed -i.tmp "s|^AUTH_MAX_NB_USER=.*|AUTH_MAX_NB_USER='${max_users}'|" "$env_file"
                else
                    # Replace without quotes
                    sed -i.tmp "s|^AUTH_MAX_NB_USER=.*|AUTH_MAX_NB_USER=${max_users}|" "$env_file"
                fi
                echo -e "${GREEN}Updated AUTH_MAX_NB_USER to ${max_users}${NC}"
            else
                # Add new variable with quotes (default format)
                echo "AUTH_MAX_NB_USER=\"${max_users}\"" >> "$env_file"
                echo -e "${GREEN}Added AUTH_MAX_NB_USER=\"${max_users}\"${NC}"
            fi
            
                # Remove temporary file created by sed
                rm -f "${env_file}.tmp"
            else
                echo -e "${RED}Error: Please enter a valid number. Skipping AUTH_MAX_NB_USER configuration.${NC}"
            fi
        else
            echo -e "${RED}Error: Maximum number of users cannot be empty.${NC}"
        fi
    else
        echo -e "${YELLOW}Skipped AUTH_MAX_NB_USER configuration.${NC}"
    fi
}
