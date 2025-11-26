package config

import (
	"os"
	"strconv"
)

func GetThresholdAlertQty() int {
	if val := os.Getenv("THRESHOLD_ALERT_QTY"); val != "" {
		if qty, err := strconv.Atoi(val); err == nil {
			return qty
		}
	}
	return 10
}

func GetJWTSecretKey() string {
	if val := os.Getenv("JWT_SECRET_KEY"); val != "" {
		return val
	}
	return "myJWTsecret@1234"
}

func GetAdminEmail() string {
	if val := os.Getenv("ADMIN_EMAIL"); val != "" {
		return val
	}
	return "aman@wg.com"
}

func GetAdminPassword() string {
	if val := os.Getenv("ADMIN_PASSWORD"); val != "" {
		return val
	}
	return "1234567"
}

func GetProductsTableName() string {
	if val := os.Getenv("PRODUCTS_TABLE_NAME"); val != "" {
		return val
	}
	return "Products"
}

func GetManagersTableName() string {
	if val := os.Getenv("MANAGERS_TABLE_NAME"); val != "" {
		return val
	}
	return "Managers"
}

func GetAWSRegion() string {
	if val := os.Getenv("AWS_REGION"); val != "" {
		return val
	}
	return "us-east-1"
}
