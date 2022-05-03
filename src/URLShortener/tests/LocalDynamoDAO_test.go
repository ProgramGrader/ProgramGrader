package tests

import "testing"

const TABLENAME = "S3URLS"

func TestCreateLocalClient(t *testing.T) {
	_, err := CreateLocalClient()
	if err != nil {
		t.Errorf("TestCreateLocalClient(), Failed. Error creating local client. got %s\n", err)
	}
}

// make delete table function move does table exist here and test using those two
func TestCreateTable(t *testing.T) {
	config, _ := CreateLocalClient()
	CreateTable(config, "test")
	if !TableExists(config, "test") {
		t.Fatalf("TestCreateTable(), Failed. Expected %s to exist.")
	}
}

func TestPutAndGet(t *testing.T) {

	type putItem struct {
		key      string
		value    string
		expected string
	}

	putItems := []putItem{
		{"key", "url", "url"},
		{"key", "update_url", "update_url"},
		{"asefay", "update_url", "update_url"},
	}

	for _, test := range putItems {
		Put(TABLENAME, test.key, test.value)
		storedItem := Get(TABLENAME, test.key)
		if storedItem != test.expected {
			t.Fatalf("TestPut(), Failed. Expected value was not found. Got %s expected %s", test.value, test.expected)
		}
	}

	DeleteAll(TABLENAME)
}

func TestDelete(t *testing.T) {
	type deleteItem struct {
		key      string
		value    string
		expected error
	}

	deleteItems := []deleteItem{
		{"key", "url", nil},
	}

	for _, test := range deleteItems {
		Put(TABLENAME, test.key, test.value)
		deleteErr := Delete(TABLENAME, test.key)
		if deleteErr != nil {
			t.Fatalf("TestDelete(), Failed. Expected error to be nil. Got %v", deleteErr)
		}
	}
}
