package src

type User struct {
	Id   string `json:"ID,omitempty" dynamodbav:"ID, omitempty"`
	Name string `json:"name,omitempty" dynamodbav:"name, omitemptye"`
	Age  string `json:"age,omitempty" dynamodbav:"age, omitempty"`
}
