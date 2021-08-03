package dockerhub

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
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

func (client *client) DeleteImages(organization string, repository string, digests []string, timeBefore time.Time, dryRun bool) {
	var manifests []*Manifest

	for _, d := range digests {
		manifests = append(manifests, &Manifest{
			Repository: repository,
			Digest:     d,
		})
	}

	post := &DeleteImagesRequest{
		DryRun:     dryRun,
		ActiveFrom: timeBefore,
		Manifests:  manifests,
	}

	body, err := json.Marshal(post)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(body))

	req, err := http.NewRequest("POST", "https://hub.docker.com/v2/namespaces/"+organization+"/delete-images", bytes.NewReader(body))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Authorization", "JWT "+client.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	rsp, _ := ioutil.ReadAll(resp.Body)

	fmt.Println(string(rsp))

	if resp.StatusCode != http.StatusOK {
		panic(string(rsp))
	}
}

func (client *client) GetImages(organization string, repository string, page int, timeBefore time.Time) ImageList {
	pageString := strconv.Itoa(page)
	timeFrom := url.QueryEscape(timeBefore.Format(time.RFC3339))

	req, err := http.NewRequest("GET", "https://hub.docker.com/v2/namespaces/"+organization+"/repositories/"+repository+"/images?page="+pageString+"&page_size=10&currently_tagged=false&ordering=last_activity&status=inactive&active_from="+timeFrom, nil)
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

	if resp.StatusCode != http.StatusOK {
		panic(string(rsp))
	}

	var tagsList ImageList
	json.Unmarshal(rsp, &tagsList)

	return tagsList
}
