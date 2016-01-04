package auth

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"payup/database"
	"payup/models"
)

// FacebookUserInfo represents the model of a user that is returned
// from facebook's oauth
type FacebookUserInfo struct {
	UserID      string `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	Gender      string `json:"gender"`
	AvatarURL   string
	AccessToken string
}

type facebookPicture struct {
	PicData data `json:"data"`
}

type data struct {
	URL string `json:"url"`
}

// GetFacebookAccessToken takes an authCode and a referrer to get the accessToken
func GetFacebookAccessToken(authCode string, referrer string) (string, error) {
	v := url.Values{}
	v.Set("client_id", os.Getenv("FACEBOOK_KEY"))
	v.Set("client_secret", os.Getenv("FACEBOOK_SECRET"))
	v.Set("redirect_uri", referrer)
	v.Set("code", authCode)

	accessTokenURL := "https://graph.facebook.com/oauth/access_token?" + v.Encode()

	resp, _ := http.Get(accessTokenURL)

	defer resp.Body.Close()
	contents, _ := ioutil.ReadAll(resp.Body)
	m, _ := url.ParseQuery(string(contents))

	if resp.StatusCode != 200 || m["access_token"][0] == "" {
		return "", errors.New("Failed to get AccessToken")
	}

	return m["access_token"][0], nil
}

// GetFacebookUserInfo validates an authCode that
// is sent from a client
func GetFacebookUserInfo(accessToken string) (FacebookUserInfo, error) {
	if accessToken == "" {
		return FacebookUserInfo{}, errors.New("AccessToken is empty")
	}

	s := url.Values{}
	s.Set("fields", "id,name,email,gender")
	s.Set("access_token", accessToken)

	resp, _ := http.Get("https://graph.facebook.com/me?" + s.Encode())

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	var userInfo FacebookUserInfo
	json.Unmarshal(body, &userInfo)

	// Get AvatarURL
	picResp, _ := http.Get("https://graph.facebook.com/" + userInfo.UserID + "/picture?type=large&redirect=false")

	defer picResp.Body.Close()
	picBody, _ := ioutil.ReadAll(picResp.Body)

	var fbPic facebookPicture
	json.Unmarshal(picBody, &fbPic)

	// Set userInfo AvatarURL
	if fbPic.PicData.URL != "" {
		userInfo.AvatarURL = fbPic.PicData.URL
	}

	return userInfo, nil
}

// UpdateFacebookUserFriends takes an authCode and an already created user, and
// finds all their facebook friends and adds them into the the database.
func UpdateFacebookUserFriends(fbAccessToken string, user models.User) {
	s := url.Values{}
	s.Set("access_token", fbAccessToken)

	res, _ := http.Get("https://graph.facebook.com/me/friends?" + s.Encode())

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	var v map[string]interface{}
	json.Unmarshal(body, &v)
	var friends []interface{}
	friends = v["data"].([]interface{})

	// TODO: Optimize this, it kinda sucks.
	for _, friend := range friends {
		friendExtID := friend.(map[string]interface{})["id"].(string)
		var friendDB models.User
		var friendship models.Friendship
		if !database.DBCon.Where("external_id = ?", friendExtID).First(&friendDB).RecordNotFound() {
			database.DBCon.Where(models.Friendship{
				UserID:   user.ID,
				FriendID: friendDB.ID,
			}).FirstOrCreate(&friendship)
		}
	}
}
