package main

import "net/http"
import "errors"

// ParseRequestForSession returns a SessionTemplateData object
// Session is not found:  returns (nil, nil)
// Session is found:  returns (session, nil)
// An error occured: returns (nil, err)
func ParseRequestForSession(req *http.Request) (*SessionTemplateData, error) {
	session := ParseSession(req)

	if session == nil {
		return nil, errors.New("Session Not Found")
	}

	return GetSessionTemplateData(*session)
}

//ParseRequestForUseData returns an User object
func ParseRequestForUseData(req *http.Request) (*User, error) {
	sessionID := ParseSession(req)

	if sessionID == nil {
		return nil, errors.New("Unknown user")
	}

	return GetSessionUserData(*sessionID)
}

// ParseSession gets session from Cookie
func ParseSession(req *http.Request) *string {
	cookies := req.Cookies()
	var session *string

	for idx := range cookies {
		if cookies[idx].Name == HttpCookieSession {
			session = &cookies[idx].Value
		}
	}

	return session
}
