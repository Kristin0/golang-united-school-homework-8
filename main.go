package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

type Arguments map[string]string

type User struct {
	Id    string `json:"id"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

type Users []User

func parseArgs() Arguments {
	var args = Arguments{}
	var p = map[string]*string{}

	p["item"] = flag.String("item", "", "")
	p["operation"] = flag.String("operation", "", "add list findById remove")
	p["fileName"] = flag.String("fileName", "", "")
	p["id"] = flag.String("id", "", "")

	flag.Parse()

	for key, value := range p {
		args[key] = *value
	}

	return args
}

func Perform(args Arguments, writer io.Writer) error {
	fileName := args["fileName"]
	operation := args["operation"]
	item := args["item"]
	id := args["id"]

	if fileName == "" {
		return errors.New("-fileName flag has to be specified")
	}

	if operation == "" {
		return errors.New("-operation flag has to be specified")
	}

	if operation != "list" &&
		operation != "add" &&
		operation != "remove" &&
		operation != "findById" {
		return errors.New("Operation " + operation + " not allowed!")

	}

	if operation == "add" {
		if item == "" {
			return errors.New("-item flag has to be specified")
		}
	}

	if operation == "remove" || operation == "findById" {
		if id == "" {
			return errors.New("-id flag has to be specified")
		}
	}

	usersFile, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0755)

	if err != nil {
		panic(err)
	}

	defer usersFile.Close()

	usersDataString, err := ioutil.ReadAll(usersFile)

	if err != nil {
		panic(err)
	}

	var users = Users{}

	if len(usersDataString) > 0 {
		err = json.Unmarshal(usersDataString, &users)
		if err != nil {
			panic(err)
		}
	}

	switch operation {
	case "list":
		result, err := json.Marshal(users)

		if err != nil {
			panic(err)
		}

		fmt.Fprint(writer, string(result))

	case "add":
		itemUser := User{}
		err := json.Unmarshal([]byte(item), &itemUser)
		if err != nil {
			panic(err)
		}

		for _, usr := range users {

			if usr.Id == itemUser.Id {
				errstr := "Item with id " + itemUser.Id + " already exists"
				writer.Write([]byte(errstr))
				return nil
			}
		}

		users = append(users, itemUser)

		datad, err := json.Marshal(users)
		if err != nil {
			panic(err)
		}

		usersFile.Write(datad)

		return nil
	case "findById":

		for _, user := range users {
			if id == user.Id {
				res, er := json.Marshal(user)
				if er != nil {
					panic(er)
				}
				writer.Write(res)
				break
			}
		}
	case "remove":
		res := Users{}
		found := false
		current := 0
		for _, user := range users {

			if id == user.Id {
				found = true
				continue
			}

			res = append(res, user)
			current++

		}
		if !found {
			writer.Write([]byte("Item with id " + id + " not found"))
			return nil
		}
		dt, err := json.Marshal(res)

		if err != nil {
			panic(err)
		}

		usersFile.Truncate(0)
		usersFile.Seek(0, 0)
		usersFile.Write(dt)
		usersFile.Close()

	}

	return nil
}

func main() {
	err := Perform(parseArgs(), os.Stdout)
	if err != nil {
		panic(err)
	}
}
