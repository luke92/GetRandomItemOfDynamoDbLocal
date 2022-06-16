# GetRandomItemOfDynamoDbLocal
Get a random record from dynamo db local

# Create Project
- Run `go mod init github.com/luke92/GetRandomItemOfDynamoDbLocal`
- Run `go get .`
- Run `go mod tidy`

# Run DynamoDBLocal

## If you're using Windows PowerShell, be sure to enclose the parameter name or the entire name and value like this:

- Run in a Powershell (or run `ryndynamo.ps1`)
- `java -D"java.library.path=./DynamoDBLocal_lib" -jar DynamoDBLocal.jar -sharedDb`

## Or you can use Docker
```shell
docker run -p 8000:8000 amazon/dynamodb-local
```