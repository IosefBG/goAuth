#!/bin/bash

# Set debug mode (true or false)
debug=false

# Start timer
start=$(date +%s)

# Function to execute command with optional debug output
execute_command() {
    local command="$1"
    local command_name="$2"

    echo "Started executing $command_name..."
    if $debug; then
        $command
    else
        $command > /dev/null 2>&1
    fi
    echo "Finished executing $command_name."
}

# Stop and remove all containers and volumes
execute_command "docker-compose down --volumes" "Remove existing containers and volumes"

# Build the Docker image
execute_command "docker build -t go-auth ./" "Docker image build"

# Check if the Docker image has been built
while true; do
    if docker images | grep -q "go-auth"; then
        echo "Docker image build completed."
        break
    else
        echo "Waiting for Docker image build to complete..."
        sleep 5  # Wait for 5 seconds before checking again
    fi
done

# Run Docker Compose
execute_command "docker compose --file .\docker-compose.yaml --project-name goauth up --detach" "Docker Compose"

# Wait for the user to press any key if debug mode is disabled

# Calculate elapsed time
end=$(date +%s)
runtime=$((end-start))

# Print elapsed time
echo "Total time elapsed: $runtime seconds."
if ! $debug; then
    read -n 1 -s -r -p "Press any key to exit..."
fi