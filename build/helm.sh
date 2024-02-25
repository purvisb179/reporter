#!/bin/bash

# Check if we're in the build directory and warn the user if they are
current_directory=$(basename "$PWD")
if [ "$current_directory" = "build" ]; then
    echo "This script should not be run from the build directory. Please move to the root directory of the project."
    exit 1
fi

# Always prompt for Helm release name and Kubernetes namespace
read -p "Helm release name (e.g., myapp-release): " HELM_RELEASE_NAME
read -p "Kubernetes namespace (e.g., default): " K8S_NAMESPACE

# Define the path to the configuration file
CONFIG_FILE="./local/helm.yaml"

# Check if the configuration file exists
if [ -f "$CONFIG_FILE" ]; then
    echo "Configuration file found: $CONFIG_FILE"
    VALUES_ARG="-f $CONFIG_FILE"
else
    echo "No configuration file found. Proceeding with prompts for deployment configuration."

    # Prompt for additional deployment configuration
    read -p "Docker image repository (e.g., username/repository): " IMAGE_REPOSITORY
    read -p "Docker image tag (e.g., latest): " IMAGE_TAG
    echo "Optional: Set additional values. Leave blank if none."
    read -p "Set values (e.g., key1=value1,key2=value2): " SET_VALUES

    # Convert comma-separated key=value pairs into --set arguments for Helm
    IFS=',' read -ra ADDR <<< "$SET_VALUES"
    for i in "${ADDR[@]}"; do
        SET_ARGS+=" --set $i"
    done

    VALUES_ARG="$SET_ARGS"
fi

# Confirm before proceeding
echo "You're about to deploy with the following configuration:"
if [ -f "$CONFIG_FILE" ]; then
    echo "Using configuration from $CONFIG_FILE"
else
    echo "Helm Release Name: $HELM_RELEASE_NAME"
    echo "Namespace: $K8S_NAMESPACE"
    echo "Image: $IMAGE_REPOSITORY:$IMAGE_TAG"
    if [ -n "$SET_VALUES" ]; then
        echo "Additional Values: $SET_VALUES"
    fi
fi
read -p "Proceed? (y/N): " confirm && [[ $confirm == [yY] ]] || exit 1

# Deploy with Helm
echo "Deploying with Helm..."
helm upgrade --install "$HELM_RELEASE_NAME" ./helm \
    --namespace "$K8S_NAMESPACE" \
    $VALUES_ARG \
    --wait

if [ $? -ne 0 ]; then
    echo "Helm deploy failed."
    exit 1
fi

echo "Helm deploy succeeded."
