package dockerhub

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
)

type client struct {
	token string
}

func NewClient(auth Auth) *client {
	c := &client{}
	c.authorize(auth)

	return c
}

func (client *client) authorize(auth Auth) {
	payloadBytes, err := json.Marshal(auth)
	if err != nil {
		panic(err)
	}
	body := bytes.NewReader(payloadBytes)

	req, err := http.NewRequest("POST", "https://hub.docker.com/v2/users/login/", body)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var token Token
	plainBody, _ := ioutil.ReadAll(resp.Body)

	json.Unmarshal(plainBody, &token)

	client.token = token.Token
}

func (client *client) DeleteTag(organization string, repository string, tag string) {
	req, err := http.NewRequest("DELETE", "https://hub.docker.com/v2/repositories/"+organization+"/"+repository+"/tags/"+tag+"/", nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Authorization", "JWT "+client.token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
}

func (client *client) GetTags(organization string, repository string, page int) TagsList {
	pageString := strconv.Itoa(page)

	req, err := http.NewRequest("GET", "https://hub.docker.com/v2/repositories/"+organization+"/"+repository+"/tags/?page="+pageString, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Authorization", "JWT "+client.token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	rsp, _ := ioutil.ReadAll(resp.Body)

	var tagsList TagsList
	json.Unmarshal(rsp, &tagsList)

	return tagsList
}
