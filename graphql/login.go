package graphql

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/graphql-go/graphql"
)

func (s *Server) loginMutator(params graphql.ResolveParams) (interface{}, error) {
	logger.Debugf("trying to login")

	// marshall and cast the argument values
	username, _ := params.Args["username"].(string)
	password, _ := params.Args["password"].(string)

	user, err := s.db.ReadUserByUsername(username)
	if err != nil {
		logger.Errorf("db error: %s", err.Error())
		return nil, err
	}
	if user == nil {
		return nil, ErrBadLogin
	}

	// check password validity
	if !user.CheckPasswordHash(password) {
		return nil, ErrBadLogin
	}

	// create jwt
	ts, err := s.createToken(user)
	if err != nil {
		logger.Debugf("error creating token: %s", err.Error())
		return nil, err
	}

	// save jwt
	err = s.createAuth(user.ID, ts)
	if err != nil {
		logger.Debugf("error saving token: %s", err.Error())
		return nil, err
	}

	return ts, nil
}

func (s *Server) logoutMutator(params graphql.ResolveParams) (interface{}, error) {
	logger.Debugf("trying to logout")

	if params.Context.Value(MetadataKey) == nil {
		return nil, ErrUnauthorized
	}
	metadata := params.Context.Value(MetadataKey).(*AccessDetails)

	err := s.deleteTokens(metadata)
	if err != nil {
		logger.Tracef("can't delete tokens: %s", err.Error())
		return nil, err
	}

	result := map[string]bool{"success": true}

	return &result, nil
}

func (s *Server) refreshAccessTokenMutator(params graphql.ResolveParams) (interface{}, error) {
	logger.Debugf("trying to refresh token")

	// marshall and cast the argument values
	refreshToken, _ := params.Args["refreshToken"].(string)

	//verify the token
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.refreshSecret, nil
	})

	//if there is an error, the token must have expired
	if err != nil {
		logger.Tracef("token error: %s", err.Error())
		return nil, err
	}
	//is token valid?
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return nil, ErrUnauthorized
	}

	//Since token is valid, get the uuid:
	claims, ok := token.Claims.(jwt.MapClaims) //the token claims should conform to MapClaims
	if ok && token.Valid {
		// read key data
		refreshString, ok := claims[ClaimRefreshId].(string) //convert the interface to string
		if !ok {
			logger.Tracef("claim %s missing", ClaimRefreshId)
			return nil, ErrUnprocessableEntity
		}
		userId, err := uuid.Parse(claims[ClaimUserId].(string))
		if err != nil {
			logger.Tracef("%s is an invalid uuid: %s", claims[ClaimUserId].(string), err.Error())
			return nil, err
		}

		// get user
		user, err := s.db.ReadUser(userId)
		if err != nil {
			logger.Errorf("getting user: %s", err.Error())
			return nil, err
		}
		if user == nil {
			return nil, ErrUnauthorized
		}

		// Delete the previous Refresh Token
		deleted, err := s.redis.DelRefreshToken(refreshString)
		if err != nil {
			logger.Errorf("redis error: %s", err.Error())
			return nil, err
		}
		if deleted == 0 {
			msg := "something went wrong"
			logger.Tracef(msg)
			return nil, errors.New(msg)
		}

		// Create new pairs of refresh and access tokens
		ts, createErr := s.createToken(user)
		if createErr != nil {
			logger.Tracef("error creating token: %s", createErr)
			return nil, createErr
		}

		// save the tokens metadata to redis
		saveErr := s.createAuth(userId, ts)
		if saveErr != nil {
			logger.Tracef("error saving token: %s", createErr)
			return nil, saveErr
		}

		return ts, nil
	}

	return nil, ErrRefreshExpired
}