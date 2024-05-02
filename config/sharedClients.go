package config

// All dummy code, just to prove viability of sharing single client instances
type SharedClients struct {
	DynamoDbclient   *DynamoDBClient
	OpensearchClient *OpensearchClient
	S3Client         *S3Client
	KcAdminClient    *KcAdminClient
}

type DynamoDBClient struct {
	port string
}

type OpensearchClient struct {
	url string
}

type S3Client struct {
	region string
}

type KcAdminClient struct {
	url string
}

// Define funcs to create clients using env vars here
