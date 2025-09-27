#!/bin/bash

# -e: Exit immediately if a command exits with a non-zero status.
# -o pipefail: The return value of a pipeline is the status of the last command
# to exit with a non-zero status, or zero if no command exited with a non-zero status.
set -eo pipefail

# --- Configuration ---
# Use default value "up" if no direction is provided.
DIRECTION=${1:-up}
MIGRATIONS_DIR="/app/migrations"
RESTRICTED_WORDS_FILE="/app/restricted_words.txt"

# --- Database Migrations ---
echo "Running migrations in direction: $DIRECTION"
# Quote DB_URL to handle potential special characters in the connection string.
goose -dir "$MIGRATIONS_DIR" postgres "$DB_URL" "$DIRECTION"

# --- Seed Restricted Words ---
# First, check if the restricted words file exists and is readable.
if [[ -f "$RESTRICTED_WORDS_FILE" ]]; then
  echo "Seeding restricted words from $RESTRICTED_WORDS_FILE..."

  # Generate a single SQL script to insert all words in one transaction.
  # This is far more efficient than connecting for each word.
  # We use a subshell and process substitution to build the script.
  psql "$DB_URL" <<EOF
BEGIN;
$(while IFS= read -r word || [[ -n "$word" ]]; do
    # 1. Remove carriage returns that can be left by some editors
    word=$(echo "$word" | tr -d '\r')
    # 2. Skip empty lines
    if [[ -z "$word" ]]; then
      continue
    fi
    # 3. Escape single quotes to prevent SQL syntax errors and injection.
    sanitized_word=$(echo "$word" | sed "s/'/''/g")
    # 4. Generate the INSERT statement.
    echo "INSERT INTO restricted_words (word) VALUES ('$sanitized_word') ON CONFLICT (word) DO NOTHING;"
done < "$RESTRICTED_WORDS_FILE")
COMMIT;
EOF

  echo "Finished seeding restricted words."
else
  echo "Warning: Restricted words file not found. Skipping seed."
fi

# --- Start Application ---
echo "Starting application..."
# Use exec to replace the shell process with the server process.
# This is a best practice for container entrypoints.
exec /bin/server
