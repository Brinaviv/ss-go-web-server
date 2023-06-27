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
	relativePath = "/users"
	idParamName  = "id"
)

func (ctrl *UsersController) Register(router *gin.RouterGroup) {
	group := router.Group(relativePath)
	group.GET(fmt.Sprintf("/:%s", idParamName), ctrl.getUserByIdHandler)
	group.POST("", ctrl.createUserHandler)
	group.PATCH(fmt.Sprintf("/:%s", idParamName), ctrl.updateUserHandler)
	group.DELETE(fmt.Sprintf("/:%s", idParamName), ctrl.deleteUserHandler)
}

func (ctrl *UsersController) getUserByIdHandler(ctx *gin.Context) {
	id, err := getUserIDParam(ctx)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	userEntity, err := ctrl.UserDAO.FetchByID(id)
	if err != nil {
		var errNotFound *dal.ErrIDNotFound
		if errors.As(err, &errNotFound) {
			ctx.AbortWithError(http.StatusNotFound, errNotFound)
		} else {
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

	id, err := getUserIDParam(ctx)
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
	id, err := getUserIDParam(ctx)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	err = ctrl.UserDAO.Delete(id)
	if err != nil {
		var errNotFound *dal.ErrIDNotFound
		if errors.As(err, &errNotFound) {
			ctx.AbortWithError(http.StatusNotFound, errNotFound)
		} else {
			ctx.AbortWithError(http.StatusInternalServerError, err)
		}
		return
	}
	ctx.IndentedJSON(http.StatusOK, &struct {
		Msg string `json:"message"`
	}{Msg: fmt.Sprintf("deleted user with id %s", id.String())})
}

func getUserIDParam(ctx *gin.Context) (model.UserID, error) {
	return model.ParseUserID(ctx.Param(idParamName))
}
