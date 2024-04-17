#!/bin/bash

# Placeholder for fetching public key from metadata and checking AWS Secrets Manager
PUBLIC_KEY=$(curl http://169.254.170.2/v2/metadata | jq -r '.Containers[0].Labels.public_key')
if [ -n "$PUBLIC_KEY" ]; then
    aws secretsmanager get-secret-value --secret-id $PUBLIC_KEY --query 'SecretString' --output text > /home/masa/.masa/masa_oracle_key
fi

# Start your main application
exec /usr/bin/masa-node --bootnodes="$BOOTNODES" --env="$ENV" --writerNode="$WRITER_NODE" --cachePath="$CACHE_PATH"

