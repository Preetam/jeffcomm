package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/smtp"
	"regexp"
	"strings"
)

func parseBalance(line string) float32 {
	reg := regexp.MustCompile(`-?[$][\d,]+[.]\d\d`)
	str := reg.FindString(line)

	invalidReg := regexp.MustCompile(`[$,]`)
	str = invalidReg.ReplaceAllString(str, "")

	f := float32(0)
	fmt.Sscanf(str, "%f", &f)
	return f
}

func main() {
	username := flag.String("username", "", "Jefferson Commons portal username")
	password := flag.String("password", "", "Jefferson Commons portal password")

	smtpAddr := flag.String("smtp-addr", "", "Address of the SMTP server")
	smtpHost := flag.String("smtp-host", "", "")
	smtpUsername := flag.String("smtp-username", "", "")
	smtpPassword := flag.String("smtp-password", "", "")
	sender := flag.String("sender", "", "")
	mailTo := flag.String("mail-to", "", "")

	forceEmail := flag.Bool("force-send-email", false, "")
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

	res, err = client.Get("https://jeffersoncommons.residentportal.com/resident_portal/?module=home&action=show_resident_alert_balance&lease_status[type_id]=4")
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
			f := parseBalance(line)

			fmt.Printf("\nOur balance is $%.2f\n", f)

			if f <= 0 && !*forceEmail {
				break
			}

			smtpAuth := smtp.PlainAuth("", *smtpUsername, *smtpPassword, *smtpHost)
			to := strings.Split(*mailTo, ",")

			log.Println("sending mail to ", to, "with the following auth ", smtpAuth)

			err := smtp.SendMail(*smtpAddr, smtpAuth, *sender, to, []byte(fmt.Sprintf(`From: %s
To: %s
Subject: Jefferson Commons Balance

Hello! Our Jefferson Commons balance is $%.2f. Remember to pay it on time!
`, *sender, *mailTo, f)))

			if err != nil {
				log.Fatal(err)
			}

			break
		}
	}

	res.Body.Close()
}
