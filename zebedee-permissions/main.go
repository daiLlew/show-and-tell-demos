package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	. "github.com/logrusorgru/aurora"
	"github.com/pkg/errors"
)

const (
	datasetID    = "suicides-in-the-uk"
	collectionID = "suicides-7d357df1521479c2d01a324e1ad78702c25308d88be2eb0a4bbf520122b7f42e"
)

var (
	c http.Client
	w *tabwriter.Writer
)

func init() {
	c = http.Client{Timeout: time.Second * 5}
	w = tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.TabIndent|tabwriter.Debug)
}

type Profile struct {
	Email       string
	Name        string
	Permissions CRUD
}

type CRUD struct {
	Permissions []string `json:"permissions"`
}

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func main() {
	profiles := []*Profile{
		{
			Email: "d.llewellyn@ons.gov.uk",
			Name:  "Publisher",
		},
		{
			Email: "a@ons.gov.uk",
			Name:  "Viewer",
		},
		{
			Email: "b@ons.gov.uk",
			Name:  "Viewer",
		},
	}

	fmt.Println()
	fmt.Fprintf(w, " %s\t %s\t %s\t %s\t %s\n", Col1("User Type"), Col2("Email"), Col1("Collection"), Col2("Dataset"), Col1("Permissions Granted"))

	for _, p := range profiles {
		session := login(p.Email)
		p.getPermissions(session)
		fmt.Fprintf(w, " %s\t %s\t %s\t %s\t %s\n", Col1(p.Name), Col2(p.Email), Col1(collectionID), Col2(datasetID), Col1(strings.Join(p.Permissions.Permissions, ", ")))
	}
	fmt.Fprint(w, "\n")
	w.Flush()
}

func (p *Profile) getPermissions(session string) {
	url := fmt.Sprintf("http://localhost:8082/permissions?dataset_id=%s&collection_id=%s", datasetID, collectionID)
	r, err := http.NewRequest("GET", url, nil)
	if err != nil {
		errorExit(err)
	}

	r.Header.Set("X-Florence-Token", session)

	resp, err := c.Do(r)
	if err != nil {
		errorExit(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		errorExit(err)
	}

	var perms CRUD
	err = json.Unmarshal(body, &perms)
	if err != nil {
		errorExit(err)
	}

	p.Permissions = perms
}

func login(email string) string {
	body := Credentials{Email: email, Password: "one two three four"}

	b, err := json.Marshal(body)
	if err != nil {
		errorExit(err)
	}

	r, err := http.NewRequest("POST", "http://localhost:8082/login", bytes.NewReader(b))
	if err != nil {
		errorExit(err)
	}

	resp, err := c.Do(r)
	if err != nil {
		errorExit(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		errorExit(errors.New("non 200 status for login"))
	}

	rB, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		errorExit(err)
	}
	s := string(rB)
	return strings.Replace(s, "\"", "", -1)
}

func Title1(s string) string {
	return Bold(Cyan(s)).String()
}

func Col1(s string) string {
	return Cyan(s).String()
}

func Title2(s string) string {
	return Bold(Magenta(s)).String()
}

func Col2(s string) string {
	return Magenta(s).String()
}

func errorExit(err error) {
	panic(err)
}
