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
        fi
    else
        echo -e "${GREEN}All services are up to date!${NC}"
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
}
