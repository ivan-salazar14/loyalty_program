# Ecommerce plan Loyalty 
Loyalty project made with architecture DDD and single page on dynamodb.

## Steps:
1. Define your aws credentials on your local environment 
2. Create a single table with script infrastructure/loyaltyTable.
3. Start server.

### Server start 
    go run cmd/main.go
### Routes:
	GET     /loyalty/{id}

	POST    /loyalty/redeem

	POST    /loyalty/collect

    GET     /loyalty/moves