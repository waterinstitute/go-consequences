#!/bin/bash

IMAGE_NAME="thewaterinstitute/goconsequence_raster"

# Function to display usage information
show_usage() {
    echo "Usage: $0 input_path output_folder"
    echo "Options:"
    echo "  input_path      Path to a .tif file or a folder containing .tif files."
    echo "  output_folder   Path to the output folder where shapefiles will be saved."
}

# Check if the first argument is --help
if [ "$1" == "--help" ]; then
    show_usage
    exit 0
fi

# Check if two arguments are provided
if [ $# -ne 2 ]; then
    echo "Error: Invalid number of arguments."
    show_usage
    exit 1
fi

# Check if docker command is available
if ! command -v docker &> /dev/null; then
    echo "Docker is not installed. Please install Docker before running this script."
    exit 1
fi

# Check if jq is installed, and if not, install it
if ! command -v jq &> /dev/null; then
    echo "jq is not installed. Installing..."
    sudo apt update
    sudo apt install -y jq
fi

LATEST_VERSION=$(curl -s "https://registry.hub.docker.com/v2/repositories/thewaterinstitute/goconsequence_raster/tags"| jq -r '."results"[]["name"]' | sort -Vr | head -n 1)

echo "Docker image version for GoConsequence:"
echo "$LATEST_VERSION"

# Check if the latest version of the image is already downloaded
if [[ -z $(docker images -q "$IMAGE_NAME:$LATEST_VERSION") ]]; then
    echo "Latest version not found, downloading..."
    docker pull "$IMAGE_NAME:$LATEST_VERSION"
else
    echo "The latest version is already downloaded."
fi

# Check if two arguments are provided
if [ $# -ne 2 ]; then
    echo "Usage: $0 input_path(Tif file or folder) output_folder"
    exit 1
fi

INPUT="$1"
OUTPUT_FOLDER="$2"

# Function to check if a file is a .tif file
is_tif_file() {
    local file="$1"
    [[ "$file" == *.tif ]]
}

# Function to get file name without extension
get_filename_without_extension() {
    local filename
	filename=$(basename "$1")
    echo "${filename%.*}"
}
# Convert input and output paths to absolute paths
INPUT=$(readlink -f "$INPUT")
OUTPUT_FOLDER=$(readlink -f "$OUTPUT_FOLDER")

# Get Structure inventory file
structure_inventory_folder="/app/structure_inventory"
# Check if the provided input path is a file or folder
if [ -f "$INPUT" ]; then
    # Single file provided, validate that it's a .tif file
    if is_tif_file "$INPUT"; then
        tif_file="$INPUT"
        tif_filename=$(get_filename_without_extension "$tif_file")
        echo "Running command for $tif_filename"
        docker run --rm -v "$INPUT":"$INPUT" -v "$OUTPUT_FOLDER":"$OUTPUT_FOLDER" --name "goconsequence" "$IMAGE_NAME:$LATEST_VERSION" -raster "$INPUT" -si "${structure_inventory_folder}/si.shp" -result "$OUTPUT_FOLDER/${tif_filename}_si_"
    else
        echo "Provided file is not a .tif file."
        exit 1
    fi
elif [ -d "$INPUT" ]; then
    # Folder provided
    FOLDER_PATH="$INPUT"
    # Loop through .tif files in the folder and run the command for each file
    for tif_file in "$FOLDER_PATH"/*.tif; do
        if [ -f "$tif_file" ]; then
            tif_filename=$(get_filename_without_extension "$tif_file")
            echo "Running command for $tif_filename"
            docker run -v "$tif_file":"$tif_file" -v "$OUTPUT_FOLDER":"$OUTPUT_FOLDER" "$IMAGE_NAME:$LATEST_VERSION" -raster "$tif_file" -si "${structure_inventory_folder}/si.shp" -result "$OUTPUT_FOLDER/${tif_filename}_si_"
        fi
    done
else
    echo "Invalid input path. Please provide a valid .tif file or folder."
    exit 1
fi
