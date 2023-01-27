package structs

import "github.com/golang-jwt/jwt"

// struct event ที่เข้าจาก Webhook รูปแบบ Text
type LineMessage struct {
	Destination string `json:"destination"`
	Events      []struct {
		ReplyToken string `json:"replyToken"`
		Type       string `json:"type"`
		Timestamp  int64  `json:"timestamp"`
		Source     struct {
			Type   string `json:"type"`
			UserID string `json:"userId"`
		} `json:"source"`
		Message struct {
			ID   string `json:"id"`
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"message"`
	} `json:"events"`
}

// --------------------------------

// struct event ที่จะส่งไปให้ Line Replytoken
type ReplyMessage struct {
	ReplyToken string `json:"replyToken"`
	Messages   []Text `json:"messages"`
}

type Text struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type LineVerify struct {
	ClientuserIdID string `json:"client_id"`
	ExpiresIn      int64  `json:"expires_in"`
	Scope          string `json:"scope"`
}

type LineProfile struct {
	DisplayName string `json:"displayName"`
	PictureURL  string `json:"pictureUrl"`
	UserID      string `json:"userId"`
	Email       string `json:"email"`
}

type PayLoad struct {
	ISS     string   `json:"iss"`
	SUB     string   `json:"sub"`
	AUD     string   `json:"aud"`
	EXP     int64    `json:"exp"`
	IAT     int64    `json:"iat"`
	AMR     []string `json:"amr"`
	Name    string   `json:"name"`
	Picture string   `json:"picture"`
	Email   string   `json:"email"`
	jwt.Claims
}

type LineLiff struct {
	Code            string `json:"code" query:"code" form:"code"`
	State           string `json:"state" query:"state" form:"state"`
	LiffRedirectUri string `json:"liffRedirectUri" query:"liffRedirectUri" form:"liffRedirectUri"`
	LiffClientId    string `json:"liffClientId" query:"liffClientId" form:"liffClientId"`
}

///----------------------------------
