package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
	"strings"
)

func main() {
	username := flag.String("username", "", "Jefferson Commons portal username")
	password := flag.String("password", "", "Jefferson Commons portal password")
	flag.Parse()

	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatal(err)
	}
	client := http.Client{
		Jar: jar,
	}

	res, err := client.PostForm("https://jeffersoncommons.residentportal.com/resident_portal/?module=authentication&action=attempt_login&return_url=",
		map[string][]string{
			"customer[username]":              []string{*username},
			"customer[password]":              []string{*password},
			"return_url":                      []string{},
			"is_attempt_from_resident_portal": []string{"1"},
		})

	if err != nil {
		log.Fatal(err)
	}

	if res.StatusCode != 200 {
		log.Fatalf("Did not get status code 200. Got %d.", res.StatusCode)
	}

	res.Body.Close()

	res, err = client.Get("https://jeffersoncommons.residentportal.com/resident_portal/?module=home")
	if err != nil {
		log.Fatal(err)
	}

	if res.StatusCode != 200 {
		log.Fatalf("Did not get status code 200. Got %d.", res.StatusCode)
	}

	scan := bufio.NewScanner(res.Body)
	for scan.Scan() {
		line := scan.Text()
		if strings.Contains(line, "Your Balance:") {
			var f float32
			fmt.Sscanf(line[strings.Index(line, "$"):], "$%f", &f)
			fmt.Printf("\nOur balance is $%.2f\n", f)
			break
		}
	}

	res.Body.Close()
}
