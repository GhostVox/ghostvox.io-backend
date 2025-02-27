direction=$1
cd ./sql/schema

 goose postgres "postgres://localhost:5432/ghostvox?sslmode=disable" $direction
