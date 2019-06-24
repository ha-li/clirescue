package trackerapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	u "os/user"
	s "strings"
	"github.com/ha-li/clirescue/cmdutil"
	"github.com/ha-li/clirescue/user"
)

var (
	URL          string     = "https://www.pivotaltracker.com/services/v5/me"
	FileLocation string     = homeDir() + "/.tracker"
	currentUser  *user.User = user.New()
	Stdout       *os.File   = os.Stdout
	Credentials  string     = homeDir() + "/.cred"
)

func Me() {
	// maybe check for the existence of .cred file

	contents, error := ioutil.ReadFile(Credentials)

	// if error is nil means there was no errors, the file
	// exists, read the credentials
	if error == nil {
		getCredential(contents)
		//fmt.Printf( "file content %s\n", contents)
	} else {
		fmt.Println("file does not exist")
		setCredentials()
	}

	//home := homeDir()
	// fmt.Println(fmt.Sprintf("home directory is %s", home))
	// reads in the credentials from the STDIN
	//setCredentials()
	parse(makeRequest())
	ioutil.WriteFile(FileLocation, []byte(currentUser.APIToken), 0644)
}

func makeRequest() []byte {
	client := &http.Client{}
	req, err := http.NewRequest("GET", URL, nil)
	req.SetBasicAuth(currentUser.Username, currentUser.Password)
	resp, err := client.Do(req)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Print(err)
	}
	fmt.Printf("\n****\nAPI response: \n%s\n", string(body))
	return body
}

func parse(body []byte) {
	var meResp = new(MeResponse)
	err := json.Unmarshal(body, &meResp)
	if err != nil {
		fmt.Println("error:", err)
	}

	currentUser.APIToken = meResp.APIToken
}

func setCredentials() {
	fmt.Fprint(Stdout, "Username: ")
	var username = cmdutil.ReadLine()
	cmdutil.Silence()
	fmt.Fprint(Stdout, "Password: ")

	var password = cmdutil.ReadLine()

	// need to save to file

	login (username, password)
	//currentUser.Login(username, password)
	cmdutil.Unsilence()


}

func getCredential(content []byte) {
	c := string(content[:len(content)-1])
	allTokens := s.Split(c, "\n")

	uTokens := s.Split(allTokens[0], ":")
	user := uTokens[1]

	pTokens := s.Split(allTokens[1], ":")
	password := pTokens[1]

	//fmt.Printf( "user: %s, password: %s\n", user, password)
	login(user, password)
	cmdutil.Unsilence()
}

func login (username, password string) {
	currentUser.Login(username, password)
}

func homeDir() string {
	usr, _ := u.Current()
	return usr.HomeDir
}

type MeResponse struct {
	APIToken string `json:"api_token"`
	Username string `json:"username"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Initials string `json:"initials"`
	Timezone struct {
		Kind      string `json:"kind"`
		Offset    string `json:"offset"`
		OlsonName string `json:"olson_name"`
	} `json:"time_zone"`
}
