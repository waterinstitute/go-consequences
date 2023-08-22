## Executing go_consequence

## Prerequisites

- Docker installed
  
- User with permission to install or jq previously installed. The Script will run `apt install -y jq` if it doesn't find it.
  
- The script will download the environment from Docker Hub using `curl`, then, it needs internet and permissions.
  
- This script was tested on Ubuntu 20.04.6 LTS
  

## Commands

- If you want to know what are the parameters of this script you can use:
  
  ```shell-session
   ./go_consequence.sh --help
  ```
  

- If you received the error `-bash: ./go_consequence.sh: Permission denied`. Please give the execution permissions using the following command:
  
  ```shell
   chmod +x go_consequence.sh
  ```
  
- Getting the consequence results for a single tif:
  
  ```shell
  ./go_consequence.sh test.tif path/to/output/folder
  ```
  
- Getting the consequence results for a folder with tif files:
  
  ```shell
  ./go_consequence.sh input/folder/path path/to/output/folder
  ```
  

This script generates shapefiles results (.shp, .shx, .dbf, .prj) with the name of the input tif and the structure inventory version used.