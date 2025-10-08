# Setup Script Modules

This directory contains the modular components of the setup script, refactored for better maintainability and organization.

## File Structure

```
scripts/
├── README.md           # This documentation
├── utils.sh           # Common utilities and helper functions
├── download.sh        # File download functionality
├── version-check.sh   # Version checking and comparison
└── env-update.sh      # Environment file updates
```

## Module Descriptions

### utils.sh
Contains common utilities used across all modules:
- Color definitions for output formatting
- Help function and usage information
- Command line argument parsing
- Required tools validation
- GitHub token checking
- Current version retrieval from .env files

### download.sh
Handles downloading files from the GitHub repository:
- File download function with error handling
- File existence checking
- User confirmation for downloads
- Smart file mapping (docker-compose.yaml → custom names)
- .env.example → .env conversion

### version-check.sh
Manages version checking and comparison:
- GitHub Container Registry API integration
- Latest version fetching
- Image and environment variable extraction from docker-compose.yaml
- Version comparison table display
- Service status determination (up-to-date/outdated)

### env-update.sh
Handles environment file updates:
- Outdated service detection
- Update summary generation
- Automatic .env file updates with backup creation
- Quote preservation in .env files
- User confirmation for updates

## Usage

The main `setup.sh` script automatically sources all modules and orchestrates the entire process. Each module can also be used independently if needed.

## Benefits of Modular Structure

1. **Maintainability**: Each module has a single responsibility
2. **Reusability**: Functions can be used independently
3. **Testability**: Individual modules can be tested separately
4. **Readability**: Code is organized logically
5. **Extensibility**: New functionality can be added as separate modules

## Adding New Features

To add new functionality:
1. Create a new module in the `scripts/` directory
2. Source it in the main `setup.sh` script
3. Call the appropriate functions from the main execution flow
4. Update this README with the new module description
