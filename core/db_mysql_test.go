package core

import "testing"

func TestConnectMysqlDB(t *testing.T) {
	_, err := connectMysqlDB("root:root@tcp(1.117.39.176:3306)/pb_test")
	if err != nil {
		t.Fatal(err)
	}
}
