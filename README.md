# INVENTORY MANAGEMENT SYSTEM - Serverless Edition

A serverless inventory management system built with AWS Lambda, DynamoDB, and API Gateway.

## Features

1. Add products with stock quantity and price
2. Update stock when items are sold or restocked
3. Generate low-stock alerts
4. Calculate total inventory value
5. JWT-based authentication
6. Fully serverless architecture

## Architecture

- **AWS Lambda**: Serverless compute for all API endpoints
- **DynamoDB**: NoSQL database for products and managers
- **API Gateway**: REST API with custom JWT authorizer
- **AWS SAM**: Infrastructure as Code for deployment

## Prerequisites

- Go 1.24.5 or higher
- AWS CLI configured with appropriate credentials
- AWS SAM CLI installed
- Make (optional, for build automation)

## Project Structure

```
.
├── cmd/
│   └── lambda/              # Lambda function handlers
│       ├── addProduct/
│       ├── getProduct/
│       ├── getAllProducts/
│       ├── sellProduct/
│       ├── restockProduct/
│       ├── updatePrice/
│       ├── deleteProduct/
│       ├── inventoryValue/
│       ├── login/
│       └── authorizer/
├── internal/
│   ├── adapters/
│   │   ├── lambda/         # Lambda adapter layer
│   │   ├── notifier/       # Notification adapters
│   │   └── repository/     # DynamoDB repositories
│   └── core/
│       ├── domain/         # Domain models
│       ├── ports/          # Interface definitions
│       └── service/        # Business logic
├── config/                 # Configuration management
├── scripts/
│   └── seed-admin/        # Admin user seeding script
├── utils/                 # Utilities (auth, etc.)
├── template.yaml          # SAM template
└── Makefile              # Build automation

```

## Installation & Setup

### 1. Clone the Repository

```bash
git clone git@github.com:amangirdhar210/inventoryManager.git
cd inventoryManager
git checkout aws
```

### 2. Install Dependencies

```bash
go mod tidy
```

### 3. Build Lambda Functions

```bash
make build
```

This will compile all Lambda functions for Linux and create deployment packages in the `bin/` directory.

### 4. Deploy to AWS

First-time deployment (interactive):

```bash
make deploy
```

Follow the prompts to configure:

- Stack name
- AWS Region
- Confirm changes before deployment
- Allow SAM CLI IAM role creation

Subsequent deployments:

```bash
make deploy-fast
```

### 5. Seed Admin User

After deployment, initialize the database with an admin user:

```bash
make init-db
```

Default credentials:

- Email: `aman@wg.com`
- Password: `1234567`

## API Endpoints

Base URL: `https://{api-id}.execute-api.{region}.amazonaws.com/prod/`

### Authentication

#### Login

```bash
POST /login
Content-Type: application/json

{
  "email": "aman@wg.com",
  "password": "1234567"
}

Response: { "token": "jwt-token" }
```

### Product Management (Requires Authentication)

All endpoints below require `Authorization: Bearer {token}` header.

#### Add Product

```bash
POST /api/products
{
  "name": "Product Name",
  "price": 99.99,
  "quantity": 100
}
```

#### Get Product

```bash
GET /api/products/{id}
```

#### Get All Products

```bash
GET /api/products
```

#### Sell Product Units

```bash
PATCH /api/products/{id}/sell
{
  "quantity": 5
}
```

#### Restock Product

```bash
PATCH /api/products/{id}/restock
{
  "quantity": 50
}
```

#### Update Product Price

```bash
PATCH /api/products/{id}/price
{
  "price": 149.99
}
```

#### Delete Product

```bash
DELETE /api/products/{id}
```

#### Get Inventory Value

```bash
GET /api/inventory/value
```

## Environment Variables

Configure these in `template.yaml` or override during deployment:

- `PRODUCTS_TABLE_NAME`: DynamoDB products table name
- `MANAGERS_TABLE_NAME`: DynamoDB managers table name
- `JWT_SECRET_KEY`: Secret key for JWT token signing
- `THRESHOLD_ALERT_QTY`: Low stock alert threshold (default: 10)
- `ADMIN_EMAIL`: Default admin email
- `ADMIN_PASSWORD`: Default admin password
- `AWS_REGION`: AWS region (default: us-east-1)

## Development

### Run Tests

```bash
make test
```

### Validate SAM Template

```bash
make validate
```

### Local Testing

Start API locally using SAM:

```bash
make local-api
```

### Clean Build Artifacts

```bash
make clean
```

## DynamoDB Schema

### Products Table

- **Partition Key**: `id` (String)
- **Attributes**: `name`, `price`, `quantity`

### Managers Table

- **Partition Key**: `id` (String)
- **GSI**: `email-index` (email as partition key)
- **Attributes**: `email`, `password`

## Cost Optimization

This serverless architecture provides:

- **Pay-per-request pricing**: Only pay for actual API calls
- **No idle costs**: No servers running when not in use
- **Auto-scaling**: Handles traffic spikes automatically
- **DynamoDB on-demand**: Pay only for reads/writes performed

## Security Considerations

1. **JWT Authentication**: All protected endpoints require valid JWT tokens
2. **IAM Roles**: Lambda functions have minimum required permissions
3. **API Gateway Authorizer**: Custom authorizer validates tokens before invoking functions
4. **Environment Variables**: Sensitive data should use AWS Secrets Manager in production

## Production Recommendations

1. Move `JWT_SECRET_KEY` to AWS Secrets Manager
2. Use AWS CloudWatch for monitoring and alerting
3. Enable AWS X-Ray for distributed tracing
4. Set up CloudWatch Logs retention policies
5. Use AWS WAF for API Gateway protection
6. Implement rate limiting and throttling
7. Enable DynamoDB Point-in-Time Recovery
8. Set up CI/CD pipeline with GitHub Actions or AWS CodePipeline

## Troubleshooting

### Deployment Issues

Check CloudFormation stack events:

```bash
aws cloudformation describe-stack-events --stack-name {stack-name}
```

### Lambda Errors

View CloudWatch Logs:

```bash
aws logs tail /aws/lambda/{function-name} --follow
```

### DynamoDB Access

Verify IAM permissions for Lambda execution role.

## Contributing

1. Create a feature branch
2. Make your changes
3. Run tests: `make test`
4. Build: `make build`
5. Submit a pull request

## License

This project is licensed under the MIT License.
