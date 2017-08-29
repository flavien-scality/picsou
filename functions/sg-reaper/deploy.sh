#!/bin/bash
set -ex

readonly LAMBDA_NAME="sg-reaper"
# Get Virtualenv Directory Path
readonly SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
if [ -z "$VIRTUAL_ENV_DIR" ]; then
    VIRTUAL_ENV_DIR="$SCRIPT_DIR/venv"
fi

echo "Using virtualenv located in : $VIRTUAL_ENV_DIR"

# If zip artefact already exists, back it up
if [ -f $SCRIPT_DIR/$LAMBDA_NAME.zip ]; then
    mv $SCRIPT_DIR/$LAMBDA_NAME.zip $SCRIPT_DIR/$LAMBDA_NAME.zip.backup
fi

# Add virtualenv libs in new zip file
cd $VIRTUAL_ENV_DIR/lib/python3.6/site-packages
zip -r9 $SCRIPT_DIR/$LAMBDA_NAME.zip *
cd $SCRIPT_DIR

# Add python code in zip file
zip -r9 $SCRIPT_DIR/$LAMBDA_NAME.zip $LAMBDA_NAME.py

# Run terraform apply
terraform apply
