/*
* @Author: zgy
* @Date:   2023/8/14 15:33
 */
package controller

import (
	"github.com/gin-gonic/gin"
	"go_gin/forms"
	"go_gin/response"
	"go_gin/service"
	"go_gin/utils"
)

func Publish(ctx *gin.Context) {
	//参数校验
	videoForm := forms.VideoForm{}
	//首先进行获取参数，进行字段校验,
	if err := ctx.ShouldBind(&videoForm); err != nil {
		utils.HandleValidatorError(ctx, err)
		return
	}
	videoForm.Data, _ = ctx.FormFile("file")
	videoService := service.NewVideoService(ctx)
	data, msg, err := videoService.Pubish(videoForm)
	if err != nil {
		response.Err(ctx, 500, 500, msg, data)
		return
	}
	response.Success(ctx, 200, msg, data)

}
func PublishList(ctx *gin.Context) {
	//参数校验
	videoListForm := forms.VideoListForm{}
	//首先进行获取参数，进行字段校验,
	if err := ctx.ShouldBind(&videoListForm); err != nil {
		utils.HandleValidatorError(ctx, err)
		return
	}
	videoService := service.NewVideoService(ctx)
	msg, data, err := videoService.PubishList(videoListForm)
	if err != nil {
		response.Err(ctx, 500, 500, msg, data)
		return
	}
	response.Success(ctx, 200, msg, data)
}

func FavoriteAction(ctx *gin.Context) {
	//参数校验
	videoFavoriteForm := forms.VideoFavcriteForm{}
	if err := ctx.ShouldBind(&videoFavoriteForm); err != nil {
		utils.HandleValidatorError(ctx, err)
		return
	}
	videoService := service.NewVideoService(ctx)
	msg, data, err := videoService.FavoritedAction(videoFavoriteForm)
	if err != nil {
		response.Err(ctx, 500, 500, msg, data)
		return
	}
	response.Success(ctx, 200, msg, data)
}

func FavoriteList(ctx *gin.Context) {
	//参数校验
	videoFavoriteListForm := forms.VideoFavoriteListForm{}
	//首先进行获取参数，进行字段校验,
	if err := ctx.ShouldBind(&videoFavoriteListForm); err != nil {
		utils.HandleValidatorError(ctx, err)
		return
	}
	videoService := service.NewVideoService(ctx)
	msg, data, err := videoService.FavoriteListFormList(videoFavoriteListForm)
	if err != nil {
		response.Err(ctx, 500, 500, msg, data)
		return
	}
	response.Success(ctx, 200, msg, data)
}

// 视频推流接口
func Feed(ctx *gin.Context) {
	//参数校验
	feedForm := forms.FeedForm{}
	//首先进行获取参数，进行字段校验,
	if err := ctx.ShouldBind(&feedForm); err != nil {
		utils.HandleValidatorError(ctx, err)
		return
	}
	videoService := service.NewVideoService(ctx)
	msg, data, err := videoService.FeedList(feedForm)
	if err != nil {
		response.Err(ctx, 500, 500, msg, data)
		return
	}
	response.Success(ctx, 200, msg, data)
}
