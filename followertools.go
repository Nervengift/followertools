package main

import "github.com/ChimeraCoder/anaconda"
import "encoding/json"
import "io/ioutil"
import "log"
import "fmt"
import "os"
import "errors"
import "runtime"

const (
	CONSUMER_KEY    = "AAAAAAAAAAAAAAAAAAAAA"
	CONSUMER_SECRET = "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA"
)

var (
	AUTH_FILE = UserHomeDir() + string(os.PathSeparator) + ".followertools-authkey"
)

func UserHomeDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}
	return os.Getenv("HOME")
}

func saveAccessToken(accessToken string, accessTokenSecret string, filename string) {
	token := map[string]string{"accessToken": accessToken, "accessTokenSecret": accessTokenSecret}
	jsonFile, _ := json.Marshal(token)
	err := ioutil.WriteFile(filename, jsonFile, 0600)
	if err != nil {
		log.Printf("Could not save access token to %s", filename)
		log.Print(err)
	}
}

func loadAccessToken(filename string) (accessToken string, accessTokenSecret string, err error) {
	jsonFile, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Printf("Could not load access token from %s", filename)
		log.Print(err)
		return
	}
	var token map[string]string
	json.Unmarshal(jsonFile, &token)
	accessToken = token["accessToken"]
	accessTokenSecret = token["accessTokenSecret"]
	if accessToken == "" || accessTokenSecret == "" {
		err = errors.New("Could not find access token in " + filename)
		log.Printf("Could not find access token in %s", filename)
	}
	return
}

func authenticate() (accessToken string, accessTokenSecret string) {
	authURL, oauthTmp, err := anaconda.AuthorizationURL("oob") // "oob" -> pin based auth
	if err != nil {
		log.Fatal("Could not get oauth URL")
	}
	var verifier string
	fmt.Printf("Authorisation required. Please visit %s to obtain your token.\n", authURL)
	fmt.Print("Token: ")
	fmt.Scanf("%s", &verifier)
	oauthTmp, values, err := anaconda.GetCredentials(oauthTmp, verifier)
	if err != nil {
		log.Fatal("Could not finish oauth handshake.")
	}
	accessToken = values["oauth_token"][0]
	accessTokenSecret = values["oauth_token_secret"][0]
	return
}

func contains(c *anaconda.Cursor, id int64) bool {
	for _, curId := range c.Ids {
		if curId == id {
			return true
		}
	}
	return false
}

func min(x int, y int) int {
	if x > y {
		return y
	} else {
		return x
	}
}

func intersect(list1 *[]int64, list2 *[]int64) []int64 {
	hashTable := make(map[int64]bool)
	intersection := make([]int64, 0)
	for _, x := range *list1 {
		hashTable[x] = true
	}
	for _, x := range *list2 {
		if hashTable[x] {
			intersection = append(intersection, x)
		}
	}
	return intersection
}

func Follows(api *anaconda.TwitterApi, user1 anaconda.User, user2 anaconda.User) bool {
	friends, _ := api.GetFriendsUser(user1.Id, nil)
	return contains(&friends, user2.Id)
}

func CommonFriendsIds(api *anaconda.TwitterApi, user1 anaconda.User, user2 anaconda.User) []int64 {
	friends1, _ := api.GetFriendsUser(user1.Id, nil)
	friends2, _ := api.GetFriendsUser(user2.Id, nil)

	return intersect(&friends1.Ids, &friends2.Ids)
}

func CommonFriends(api *anaconda.TwitterApi, user1 anaconda.User, user2 anaconda.User) (common []anaconda.User) {
	intersection := CommonFriendsIds(api, user1, user2)
	for i := 0; i < len(intersection); i += 100 {
		users, _ := api.GetUsersLookupByIds(intersection[i:min(len(intersection), i+100)], nil)
		common = append(common, users...)
	}
	return
}

func main() {
	usage := "Usage: followertools follows|friends|connection|commonfriends|commonfriendscount|path <user1> <user2>"

	anaconda.SetConsumerKey(CONSUMER_KEY)
	anaconda.SetConsumerSecret(CONSUMER_SECRET)

	accessToken, accessTokenSecret, err := loadAccessToken(AUTH_FILE)
	if err != nil {
		accessToken, accessTokenSecret := authenticate()
		saveAccessToken(accessToken, accessTokenSecret, AUTH_FILE)
	}

	api := anaconda.NewTwitterApi(accessToken, accessTokenSecret)

	args := os.Args[1:]
	if len(args) < 3 {
		fmt.Println(usage)
		os.Exit(1)
	}

	user1, err := api.GetUsersShow(args[1], nil)
	if err != nil {
		log.Fatalf("Cannot find user %s: %s\n", args[1], err.Error())
	}

	user2, err := api.GetUsersShow(args[2], nil)
	if err != nil {
		log.Fatalf("Cannot find user %s: %s\n", args[2], err.Error())
	}

	switch args[0] {
	case "follows":
		if Follows(api, user1, user2) {
			fmt.Println("yes")
		} else {
			fmt.Println("no")
			os.Exit(1)
		}
	case "friends":
		if Follows(api, user1, user2) && Follows(api, user2, user1) {
			fmt.Println("yes")
		} else {
			fmt.Println("no")
			os.Exit(1)
		}
	case "connection":
		if Follows(api, user1, user2) {
			fmt.Printf("%s (@%s) follows %s (@%s)\n", user1.Name, user1.ScreenName, user2.Name, user2.ScreenName)
		} else {
			fmt.Printf("%s (@%s) doesn't follow %s (@%s)\n", user1.Name, user1.ScreenName, user2.Name, user2.ScreenName)
		}
		if Follows(api, user2, user1) {
			fmt.Printf("%s (@%s) follows %s (@%s)\n", user2.Name, user2.ScreenName, user1.Name, user1.ScreenName)
		} else {
			fmt.Printf("%s (@%s) doesn't follow %s (@%s)\n", user2.Name, user2.ScreenName, user1.Name, user1.ScreenName)
		}
	case "commonfriends":
		for _, user := range CommonFriends(api, user1, user2) {
			fmt.Printf("%s (@%s)\n", user.Name, user.ScreenName)
		}
	case "commonfriendscount":
		fmt.Println(len(CommonFriendsIds(api, user1, user2)))
	case "path":
		fmt.Println("not implemented yet")
		os.Exit(1)
	default:
		fmt.Println(usage)
		os.Exit(1)
	}
}
