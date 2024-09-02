#!/bin/bash

# Set debug mode (true or false)
debug=false

# Ask the user if they want to deploy the backend
read -p "Do you want to deploy the backend? (y/n): " user_input_backend

# Set deploy_backend based on user input
if [[ "$user_input_backend" == "y" || "$user_input_backend" == "Y" ]]; then
    deploy_backend=true
else
    deploy_backend=false
fi

# Ask the user if they want to deploy monitoring services
read -p "Do you want to deploy monitoring services? (y/n): " user_input_monitoring

# Set deploy_monitoring based on user input
if [[ "$user_input_monitoring" == "y" || "$user_input_monitoring" == "Y" ]]; then
    deploy_monitoring=true
else
    deploy_monitoring=false
fi

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

# Conditionally build the Docker image if deploy_backend is true
if $deploy_backend; then
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
else
    echo "Skipping Docker image build for backend."
fi

# Determine which services to deploy based on user input
services_to_deploy="postgres"

if $deploy_monitoring; then
    services_to_deploy+=" grafana prometheus cadvisor node_exporter"
fi

if $deploy_backend; then
    services_to_deploy+=" go-auth"
fi

# Run Docker Compose with selected services
execute_command "docker compose --file ./docker-compose.yaml --project-name goauth up --detach $services_to_deploy" "Docker Compose with Selected Services"

# Calculate elapsed time
end=$(date +%s)
runtime=$((end-start))

# Print elapsed time
echo "Total time elapsed: $runtime seconds."

# Wait for the user to press any key if debug mode is disabled
if ! $debug; then
    read -n 1 -s -r -p "Press any key to exit..."
fi
