package api

import (
	"errors"
	"fmt"
	"github.com/brinaviv/ss-go-web-server/pkg/dal"
	"github.com/brinaviv/ss-go-web-server/pkg/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UsersController struct {
	UserDAO dal.UserDAO
}

const (
	relativePath      = "/users"
	userIDParamName   = "id"
	targetIDParamNAme = "target"
)

var (
	ErrUnfollowSelf = errors.New("a user cannot unfollow themselves")
	ErrFollowSelf   = errors.New("a user cannot follow themselves")
)

func (ctrl *UsersController) Register(router *gin.RouterGroup) {
	group := router.Group(relativePath)
	group.POST("", ctrl.createUserHandler)

	// specific user methods
	userGroup := group.Group(fmt.Sprintf("/:%s", userIDParamName))
	userGroup.GET("", ctrl.getUserByIdHandler)
	userGroup.PATCH("", ctrl.updateUserHandler)
	userGroup.DELETE("", ctrl.deleteUserHandler)

	userGroup.POST(fmt.Sprintf("/follow/:%s", targetIDParamNAme), ctrl.followHandler)
	userGroup.POST(fmt.Sprintf("/unfollow/:%s", targetIDParamNAme), ctrl.unfollowHandler)

	// TODO: POST users/:id/tweet and save it to head of user collection. Use min stack collection (compared value is time elapsed)
	// to make post tweet O(1) and home timeline O(N) where N is the combined number of tweets to display
	// TODO: GET users/:id/timeline return tweet min stack by order
	// TODO: GET users/:id/home. Take the tweet stack of each follower and implement algorithm to merge K min stacks.
	// TODO: GET /users/popular. Go over the followers map and find the one with the biggest size. Should be able to
	// use a count map map[id]int to make it more efficient
}

func (ctrl *UsersController) getUserByIdHandler(ctx *gin.Context) {
	id, err := getIDParam(ctx)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	userEntity, err := ctrl.UserDAO.FetchByID(id)
	if err != nil {
		switch {
		case errors.Is(err, dal.ErrUserNotFound):
			ctx.AbortWithError(http.StatusNotFound, err)
		default:
			ctx.AbortWithError(http.StatusInternalServerError, err)
		}
		return
	}
	ctx.IndentedJSON(http.StatusOK, userEntity)
}

func (ctrl *UsersController) createUserHandler(ctx *gin.Context) {
	user := &model.User{}
	err := ctx.BindJSON(user)
	if err != nil {
		return
	}

	err = model.ValidateUser(user)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	entity, err := ctrl.UserDAO.Create(user)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.IndentedJSON(http.StatusCreated, entity)
}

func (ctrl *UsersController) updateUserHandler(ctx *gin.Context) {
	user := &model.User{}
	err := ctx.BindJSON(user)
	if err != nil {
		return
	}

	err = model.ValidateUser(user)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	id, err := getIDParam(ctx)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	entity, err := ctrl.UserDAO.Update(id, user)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.IndentedJSON(http.StatusOK, entity)
}

func (ctrl *UsersController) deleteUserHandler(ctx *gin.Context) {
	id, err := getIDParam(ctx)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	err = ctrl.UserDAO.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, dal.ErrUserNotFound):
			ctx.AbortWithError(http.StatusNotFound, err)
		default:
			ctx.AbortWithError(http.StatusInternalServerError, err)
		}
		return
	}
	ctx.IndentedJSON(http.StatusOK, &struct {
		Msg string `json:"message"`
	}{Msg: fmt.Sprintf("deleted user with id %s", id.String())})
}

func (ctrl *UsersController) followHandler(ctx *gin.Context) {
	userID, err := getIDParam(ctx)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	targetID, err := getTargetIDParam(ctx)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if userID == targetID {
		ctx.AbortWithError(http.StatusBadRequest, ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("%w; id %s", ErrFollowSelf, userID)))
		return
	}

	followers, err := ctrl.UserDAO.Follow(userID, targetID)
	if err != nil {
		switch {
		case errors.Is(err, dal.ErrUserNotFound):
			ctx.AbortWithError(http.StatusNotFound, err)
		case errors.Is(err, dal.ErrAlreadyFollowing):
			ctx.AbortWithError(http.StatusBadRequest, err)
		default:
			ctx.AbortWithError(http.StatusInternalServerError, err)
		}
	}

	ctx.IndentedJSON(http.StatusOK, followers)
}

func (ctrl *UsersController) unfollowHandler(ctx *gin.Context) {
	userID, err := getIDParam(ctx)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	targetID, err := getTargetIDParam(ctx)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if userID == targetID {
		ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("%w; id %s", ErrUnfollowSelf, userID))
		return
	}

	followers, err := ctrl.UserDAO.Unfollow(userID, targetID)
	if err != nil {
		switch {
		case errors.Is(err, dal.ErrUserNotFound):
			ctx.AbortWithError(http.StatusNotFound, err)
		case errors.Is(err, dal.ErrAlreadyFollowing):
			ctx.AbortWithError(http.StatusBadRequest, err)
		default:
			ctx.AbortWithError(http.StatusInternalServerError, err)
		}
	}

	ctx.IndentedJSON(http.StatusOK, followers)
}

func getIDParam(ctx *gin.Context) (model.UserID, error) {
	return getUserIDParamByName(ctx, userIDParamName)
}

func getTargetIDParam(ctx *gin.Context) (model.UserID, error) {
	return getUserIDParamByName(ctx, targetIDParamNAme)
}

func getUserIDParamByName(ctx *gin.Context, paramName string) (model.UserID, error) {
	return model.ParseUserID(ctx.Param(paramName))
}
