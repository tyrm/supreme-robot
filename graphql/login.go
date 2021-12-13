package graphql

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
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
		return nil, errBadLogin
	}

	// check password validity
	if !user.CheckPasswordHash(password) {
		return nil, errBadLogin
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

	if params.Context.Value(metadataKey) == nil {
		return nil, errUnauthorized
	}
	metadata := params.Context.Value(metadataKey).(*accessDetails)

	err := s.deleteTokens(metadata)
	if err != nil {
		logger.Tracef("can't delete tokens: %s", err.Error())
		return nil, err
	}

	return success{Success: true}, nil
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
		return nil, errUnauthorized
	}

	//Since token is valid, get the uuid:
	claims, ok := token.Claims.(jwt.MapClaims) //the token claims should conform to MapClaims
	if ok && token.Valid {
		// read key data
		refreshString, ok := claims[claimRefreshID].(string) //convert the interface to string
		if !ok {
			logger.Tracef("claim %s missing", claimRefreshID)
			return nil, errUnprocessableEntity
		}
		userID, err := uuid.Parse(claims[claimUserID].(string))
		if err != nil {
			logger.Tracef("%s is an invalid uuid: %s", claims[claimUserID].(string), err.Error())
			return nil, err
		}

		// get user
		user, err := s.db.ReadUser(userID)
		if err != nil {
			logger.Errorf("getting user: %s", err.Error())
			return nil, err
		}
		if user == nil {
			return nil, errUnauthorized
		}

		// Delete the previous Refresh Token
		deleted, err := s.kv.DeleteRefreshToken(refreshString)
		if err != nil {
			logger.Errorf("kv error: %s", err.Error())
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

		// save the tokens metadata to kv
		saveErr := s.createAuth(userID, ts)
		if saveErr != nil {
			logger.Tracef("error saving token: %s", createErr)
			return nil, saveErr
		}

		return ts, nil
	}

	return nil, errRefreshExpired
}
