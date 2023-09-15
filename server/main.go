package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
)

const jsonData = `
{
  "users": [
    {
      "id": 1,
      "name": "John Doe",
      "email": "john.doe@example.com",
      "address": {
        "street": "123 Main St",
        "city": "Anytown",
        "state": "CA",
        "postalCode": "12345"
      },
      "phone": "123-456-7890",
      "isActive": true
    },
    {
      "id": 2,
      "name": "Jane Smith",
      "email": "jane.smith@example.com",
      "address": {
        "street": "456 Elm St",
        "city": "Anycity",
        "state": "NY",
        "postalCode": "67890"
      },
      "phone": "987-654-3210",
      "isActive": false
    },
    {
      "id": 3,
      "name": "Sam Johnson",
      "email": "sam.johnson@example.com",
      "address": {
        "street": "789 Maple Ave",
        "city": "Anystate",
        "state": "TX",
        "postalCode": "11223"
      },
      "phone": "567-890-1234",
      "isActive": true
    }
  ],
  "metadata": {
    "version": "1.0",
    "timestamp": "2023-09-15T10:00:00Z",
    "totalCount": 3
  }
}
`

func insertDataHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	//var data map[string]interface{}
	//err = json.Unmarshal(body, &data)
	//if err != nil {
	//	http.Error(w, "Failed to parse JSON", http.StatusBadRequest)
	//	return
	//}

	hasher := md5.New()
	hasher.Write(body)
	md5Hash := hex.EncodeToString(hasher.Sum(nil))

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, md5Hash)
}

func fetchDataHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(jsonData))
}

func main() {
	http.HandleFunc("/insert", insertDataHandler)
	http.HandleFunc("/read", fetchDataHandler)

	fmt.Println("Starting server on :6969")
	http.ListenAndServe("127.0.0.1:6969", nil)
}
