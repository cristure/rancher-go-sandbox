package login

import (
	"time"
)

const (
	SessionCookieName = "R_SESS"
)

type Request struct {
	Description  string `json:"description"`
	ResponseType string `json:"responseType"`
	Username     string `json:"username"`
	Password     string `json:"password"`
}

type Response struct {
	AuthProvider string      `json:"authProvider"`
	BaseType     string      `json:"baseType"`
	ClusterId    interface{} `json:"clusterId"`
	Created      time.Time   `json:"created"`
	CreatedTS    int64       `json:"createdTS"`
	CreatorId    interface{} `json:"creatorId"`
	Current      bool        `json:"current"`
	Description  string      `json:"description"`
	Enabled      bool        `json:"enabled"`
	Expired      bool        `json:"expired"`
	ExpiresAt    string      `json:"expiresAt"`
	Id           string      `json:"id"`
	IsDerived    bool        `json:"isDerived"`
	Labels       struct {
		AuthnManagementCattleIoKind        string `json:"authn.management.cattle.io/kind"`
		AuthnManagementCattleIoTokenUserId string `json:"authn.management.cattle.io/token-userId"`
		CattleIoCreator                    string `json:"cattle.io/creator"`
	} `json:"labels"`
	LastUpdateTime string `json:"lastUpdateTime"`
	Links          struct {
		Self string `json:"self"`
	} `json:"links"`
	Name          string `json:"name"`
	Token         string `json:"token"`
	Ttl           int    `json:"ttl"`
	Type          string `json:"type"`
	UserId        string `json:"userId"`
	UserPrincipal string `json:"userPrincipal"`
	Uuid          string `json:"uuid"`
}
