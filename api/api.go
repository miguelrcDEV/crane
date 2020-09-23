package api

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	Types "github.com/miguelrcDEV/cmpMonit/types"
	"io/ioutil"
	"net/http"
)

var ActiveSessions []Types.MediaServerActiveSessions

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return b64.StdEncoding.EncodeToString([]byte(auth))
}

func redirectPolicyFunc(req *http.Request, via []*http.Request) error {
	req.Header.Add("Authorization", "Basic "+basicAuth("username1", "password123"))
	return nil
}

func GetActiveSessionsByServer(mediaServer Types.MediaServer) Types.MediaServerResponse {
	/* fmt.Printf("MediaServer instance: %v\n", mediaServer.Instance)
	fmt.Printf("MediaServer url: %v\n", mediaServer.Url)
	fmt.Printf("MediaServer user: %v\n", mediaServer.User)
	fmt.Printf("MediaServer secret: %v\n", mediaServer.Secret) */

	client := &http.Client{
		CheckRedirect: redirectPolicyFunc,
	}

	req, err := http.NewRequest("GET", mediaServer.Url, nil)
	req.Header.Add("Authorization", "Basic "+basicAuth(mediaServer.User, mediaServer.Secret))
	resp, err := client.Do(req)

	if err != nil {
		fmt.Printf("Error getting active sessions from MediaServer%v\n", mediaServer.Instance)
	} else {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("Error decoding response from MediaServer%v\n", mediaServer.Instance)
		} else {
			var data Types.MediaServerResponse
			err = json.Unmarshal(body, &data)
			if err != nil {
				fmt.Printf("Error unmarshalling body from MediaServer%v\n", mediaServer.Instance)
			} else {
				//fmt.Printf("Results: %v\n", data)
				data.Health = true
				return data
			}
			resp.Body.Close()
		}
	}

	var emptyData Types.MediaServerResponse
	emptyData.Health = false
	return emptyData
}

func GetActiveSessions(mediaServers []Types.MediaServer) []Types.MediaServerActiveSessions {
	var activeSessions []Types.MediaServerActiveSessions

	for i := 0; i < len(mediaServers); i++ {
		var data = GetActiveSessionsByServer(mediaServers[i])
		var mediaServerActiveSessions = Types.MediaServerActiveSessions{
			MediaServer:      mediaServers[i],
			Health:           data.Health,
			NumberOfElements: data.NumberOfElements,
			Content:          data.Content,
		}
		activeSessions = append(activeSessions, mediaServerActiveSessions)
	}

	ActiveSessions = activeSessions
	//fmt.Printf("ACTIVE SESSIONS: %v\n", ActiveSessions)

	return activeSessions
}
