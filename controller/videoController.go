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

type PublishRes struct {
	Response
}

func Publish(ctx *gin.Context) {
	//参数校验
	videoForm := forms.VideoForm{}
	////首先进行获取参数，进行字段校验,
	//if err := ctx.ShouldBind(&videoForm); err != nil {
	//	utils.HandleValidatorError(ctx, err)
	//	return
	//}
	//鉴权
	videoForm.Title = ctx.PostForm("title")
	videoForm.Data, _ = ctx.FormFile("data")
	videoService := service.NewVideoService(ctx)
	msg, _, err := videoService.Pubish(videoForm)

	if err != nil {
		res := PublishRes{Response{StatusCode: 500, StatusMsg: msg.(string)}}
		response.Err(ctx, 500, res)
		return
	}
	res := PublishRes{Response{StatusCode: 0, StatusMsg: msg.(string)}}
	response.Success(ctx, res)
}

type PublishListRes struct {
	Response
	VideoList []forms.PublishRes `json:"video_list"`
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
		res := PublishListRes{
			Response{
				StatusCode: 500,
				StatusMsg:  msg.(string),
			},
			nil,
		}
		response.Err(ctx, 500, res)
		return
	}
	res := PublishListRes{
		Response{
			StatusCode: 0,
			StatusMsg:  msg.(string),
		},
		data.([]forms.PublishRes),
	}
	response.Success(ctx, res)
}

func FavoriteAction(ctx *gin.Context) {
	//参数校验
	videoFavoriteForm := forms.VideoFavcriteForm{}
	//if err := ctx.ShouldBind(&videoFavoriteForm); err != nil {
	//	utils.HandleValidatorError(ctx, err)
	//	return
	//}
	videoFavoriteForm.ActionType = ctx.PostForm("action_type")
	videoFavoriteForm.VideoId = ctx.PostForm("video_id")
	videoService := service.NewVideoService(ctx)
	msg, _, err := videoService.FavoritedAction(videoFavoriteForm)
	if err != nil {
		res := Response{StatusCode: 500, StatusMsg: msg.(string)}
		response.Err(ctx, 500, res)
		return
	}
	res := Response{StatusCode: 0, StatusMsg: msg.(string)}
	response.Success(ctx, res)
}

type FavoriteListRes struct {
	Response
	VideoList []forms.FavoriteRes `json:"video_list"`
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
		res := FavoriteListRes{
			Response{StatusCode: 500, StatusMsg: msg.(string)},
			nil,
		}
		response.Err(ctx, 500, res)
		return
	}
	res := FavoriteListRes{
		Response{StatusCode: 500, StatusMsg: msg.(string)},
		data.([]forms.FavoriteRes),
	}
	response.Success(ctx, res)
}

type Response struct {
	StatusCode int    `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
}

type FeedRes struct {
	Response
	NextTime  int             `json:"next_time"`
	VideoList []forms.FeedRes `json:"video_list"`
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
	msg, data, time, err := videoService.FeedList(feedForm)

	if err != nil {
		res := FeedRes{
			Response{StatusCode: 500, StatusMsg: msg.(string)},
			time,
			data.([]forms.FeedRes),
		}
		response.Err(ctx, 500, res)
		return
	}
	res := FeedRes{
		Response{StatusCode: 0, StatusMsg: msg.(string)},
		time,
		data.([]forms.FeedRes),
	}
	response.Success(ctx, res)
}
