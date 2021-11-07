package graphql

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/tyrm/supreme-robot/models"
	"net/http"
	"strings"
	"time"
)

const claimGroups = "groups"
const claimRefreshID = "refresh_uuid"
const claimUserID = "user_uuid"

type accessDetails struct {
	AccessID uuid.UUID
	UserID   uuid.UUID
	Groups   []uuid.UUID
}

type tokenDetails struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	AccessUUID   uuid.UUID
	RefreshUUID  string
	AtExpires    int64
	RtExpires    int64
}

func (s *Server) createAuth(userid uuid.UUID, td *tokenDetails) error {
	at := time.Unix(td.AtExpires, 0) //converting Unix to UTC(to Time object)
	rt := time.Unix(td.RtExpires, 0)
	now := time.Now()

	errAccess := s.redis.SetAccessToken(td.AccessUUID, userid, at.Sub(now))
	if errAccess != nil {
		logger.Debugf("can't save access token: %s", errAccess.Error())
		return errAccess
	}
	errRefresh := s.redis.SetRefreshToken(td.RefreshUUID, userid, rt.Sub(now))
	if errRefresh != nil {
		logger.Debugf("can't save refresh token: %s", errRefresh.Error())
		return errRefresh
	}
	return nil
}

func (s *Server) createToken(user *models.User) (*tokenDetails, error) {
	td := &tokenDetails{}
	td.AtExpires = time.Now().Add(s.accessExpiration).Unix()
	td.AccessUUID = uuid.New()

	td.RtExpires = time.Now().Add(s.refreshExpiration).Unix()
	td.RefreshUUID = td.AccessUUID.String() + "++" + user.ID.String()

	var err error

	//Creating Access Token
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["access_uuid"] = td.AccessUUID
	atClaims[claimUserID] = user.ID
	atClaims["exp"] = td.AtExpires
	atClaims[claimGroups] = user.Groups
	at := jwt.NewWithClaims(jwt.SigningMethodHS512, atClaims)
	td.AccessToken, err = at.SignedString(s.accessSecret)
	if err != nil {
		return nil, err
	}

	//Creating Refresh Token
	rtClaims := jwt.MapClaims{}
	rtClaims[claimRefreshID] = td.RefreshUUID
	rtClaims[claimUserID] = user.ID
	rtClaims["exp"] = td.RtExpires
	rt := jwt.NewWithClaims(jwt.SigningMethodHS512, rtClaims)
	td.RefreshToken, err = rt.SignedString(s.refreshSecret)
	if err != nil {
		return nil, err
	}
	return td, nil
}

func (s *Server) deleteTokens(authD *accessDetails) error {
	//get the refresh uuid
	refreshUUID := fmt.Sprintf("%s++%s", authD.AccessID, authD.UserID)
	//delete access token
	deletedAt, err := s.redis.DeleteAccessToken(authD.AccessID)
	if err != nil {
		return err
	}
	//delete refresh token
	deletedRt, err := s.redis.DeleteRefreshToken(refreshUUID)
	if err != nil {
		return err
	}
	//When the record is deleted, the return value is 1
	logger.Tracef("deletedAt: %v, deletedRt: %v", deletedAt, deletedRt)
	if deletedAt != 1 || deletedRt != 1 {
		return errors.New("something went wrong")
	}
	return nil
}

func (s *Server) extractToken(r *http.Request) string {
	bearToken := r.Header.Get("Authorization")
	//normally Authorization the_token_xxx
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}

func (s *Server) extractTokenMetadata(r *http.Request) (*accessDetails, error) {
	token, err := s.verifyToken(r)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		accessID, err := uuid.Parse(claims["access_uuid"].(string))
		if err != nil {
			return nil, err
		}
		userID, err := uuid.Parse(claims[claimUserID].(string))
		if err != nil {
			return nil, err
		}
		groups := claims[claimGroups].([]interface{})
		groupIds := make([]uuid.UUID, len(groups))
		for i, g := range groups {
			gu, err := uuid.Parse(g.(string))
			if err != nil {
				logger.Tracef("%s is not a uuid: %s", g, err.Error())
				return nil, err
			}
			groupIds[i] = gu
		}

		return &accessDetails{
			AccessID: accessID,
			Groups:   groupIds,
			UserID:   userID,
		}, nil
	}
	return nil, err
}

func (s *Server) fetchAuth(authD *accessDetails) (uuid.UUID, error) {
	userid, err := s.redis.GetAccessToken(authD.AccessID)
	if err != nil {
		return uuid.Nil, err
	}
	if authD.UserID != userid {
		return uuid.Nil, errors.New("unauthorized")
	}
	return userid, nil
}

func (s *Server) tokenValid(r *http.Request) error {
	token, err := s.verifyToken(r)
	if err != nil {
		return err
	}
	if _, ok := token.Claims.(jwt.Claims); !ok || !token.Valid {
		return err
	}
	return nil
}

// Parse, validate, and return a token.
// keyFunc will receive the parsed token and should return the key for validating.
func (s *Server) verifyToken(r *http.Request) (*jwt.Token, error) {
	tokenString := s.extractToken(r)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.accessSecret, nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}
