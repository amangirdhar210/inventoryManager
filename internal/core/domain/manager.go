package domain

type Manager struct {
	Id       string `dynamodbav:"id"`
	Email    string `dynamodbav:"email"`
	Password string `dynamodbav:"password"`
}
