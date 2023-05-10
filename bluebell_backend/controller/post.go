package controller

import (
	"bluebell_backend/dao/mysql"
	"bluebell_backend/dao/redis"
	"bluebell_backend/logic"
	"bluebell_backend/models"
	"fmt"
	"path"
	"sort"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

// PostHandler 创建帖子
func CreatePostHandler(c *gin.Context) {
	var post models.Post
	if err := c.ShouldBindJSON(&post); err != nil {
		ResponseErrorWithMsg(c, CodeInvalidParams, err.Error())
		return
	}
	// 参数校验
	// 获取作者ID，当前请求的UserID
	userID, err := GetCurrentUserID(c)
	if err != nil {
		zap.L().Error("GetCurrentUserID() failed", zap.Error(err))
		ResponseError(c, CodeNotLogin)
		return
	}
	post.AuthorId = userID
	err = logic.CreatePost(&post)
	if err != nil {
		zap.L().Error("logic.CreatePost failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, nil)
}

// PostListHandler 帖子列表
func PostListHandler(c *gin.Context) {
	order, _ := c.GetQuery("order")
	pageStr, ok := c.GetQuery("page")
	if !ok {
		pageStr = "1"
	}
	pageNum, err := strconv.ParseInt(pageStr, 10, 64)
	if err != nil {
		pageNum = 1
	}
	posts := redis.GetPost(order, pageNum, c)
	//fmt.Println(len(posts))
	ResponseSuccess(c, posts)
}

func GetPostListBYUser(c *gin.Context) {
	page, _ := c.GetQuery("page")
	id, _ := GetCurrentUserID(c)
	postIDList := mysql.GetPostByUser(id)
	keys := make([]string, 0)
	//for _, s := range postIDList {
	//	keys = append(keys, strconv.FormatUint(s, 10))
	//}
	for i := len(postIDList) - 1; i >= 0; i-- {
		keys = append(keys, strconv.FormatUint(postIDList[i], 10))
	}
	pageNum, _ := strconv.ParseInt(page, 10, 64)
	postList := redis.GetPostBYKeys(keys, pageNum)
	ResponseSuccess(c, postList)
}

func GetOtherUserPost(c *gin.Context) {
	page, _ := c.GetQuery("page")
	id := c.Param("id")
	atom, _ := strconv.Atoi(id)
	u := uint64(atom)
	postIDList := mysql.GetPostByUser(u)
	keys := make([]string, 0)
	//for _, s := range postIDList {
	//	keys = append(keys, strconv.FormatUint(s, 10))
	//}
	for i := len(postIDList) - 1; i >= 0; i-- {
		keys = append(keys, strconv.FormatUint(postIDList[i], 10))
	}
	pageNum, _ := strconv.ParseInt(page, 10, 64)
	postList := redis.GetPostBYKeys(keys, pageNum)
	ResponseSuccess(c, postList)
}

func CommunityListHandler(c *gin.Context) {
	community, _ := c.GetQuery("community")
	order, _ := c.GetQuery("order")
	pageStr, ok := c.GetQuery("page")
	if !ok {
		pageStr = "1"
	}
	pageNum, err := strconv.ParseInt(pageStr, 10, 64)
	if err != nil {
		pageNum = 1
	}
	posts := redis.GetCommunityPost(community, order, pageNum)
	ResponseSuccess(c, posts)
}

func CommunityRankHandler(c *gin.Context) {
	communityList, _ := mysql.GetCommunityList()
	communityRank := make([]models.CommunityRankStruct, 0)
	var communityStruct models.CommunityRankStruct
	//for i := 0; i < len(communityList); i++ {
	//	communityRank[i].Community = communityList[i]
	//	communityRank[i].Rank = redis.GetCommunityNum(communityList[i].Name)
	//}
	for _, community := range communityList {
		communityStruct.Community = community
		communityStruct.Num = redis.GetCommunityNum(community.Name)
		communityRank = append(communityRank, communityStruct)
	}
	sort.Slice(communityRank, func(i, j int) bool {
		return communityRank[i].Num > communityRank[j].Num
	})
	ResponseSuccess(c, communityRank)
}

func CommunityRankToday(c *gin.Context) {
	communityList, _ := mysql.GetCommunityList()
	communityRank := make([]models.CommunityRankStruct, 0)
	var communityStruct models.CommunityRankStruct
	for _, community := range communityList {
		communityStruct.Community = community
		communityStruct.Num = redis.GetCommunityTodayNum(community.Name)
		communityRank = append(communityRank, communityStruct)
	}
	sort.Slice(communityRank, func(i, j int) bool {
		return communityRank[i].Num > communityRank[j].Num
	})
	ResponseSuccess(c, communityRank)
}

func PostList2Handler(c *gin.Context) {
	pageStr := c.Query("page")
	sizeStr := c.Query("size")

	var (
		page int64
		size int64
		err  error
	)
	page, err = strconv.ParseInt(pageStr, 10, 64)
	if err != nil {
		page = 1
	}
	size, err = strconv.ParseInt(sizeStr, 10, 64)
	if err != nil {
		size = 10
	}

	data, err := logic.GetPostList2(page, size)
	if err != nil {
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, data)

}

// PostDetailHandler 帖子详情
func PostDetailHandler(c *gin.Context) {
	postId := c.Param("id")

	post, err := logic.GetPost(postId)
	if err != nil {
		zap.L().Error("logic.GetPost(postID) failed", zap.String("postId", postId), zap.Error(err))
	}

	userId, _ := GetCurrentUserID(c)
	//查询收藏状态
	collectList := mysql.GetCollectList(userId)
	for _, collect := range collectList {
		if collect == postId {
			post.Direction = 1
		}
	}
	//查询关注状态
	concernList := mysql.GetConcernList(userId)
	for _, concern := range concernList {
		if concern == strconv.FormatUint(post.AuthorId, 10) {
			post.IsConcern = 1
		}
	}
	ResponseSuccess(c, post)
}

func UploadImg(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(200, gin.H{
			"code": 400,
			"msg":  "上传失败!",
		})
	}
	files := form.File["file"]
	for _, file := range files {
		fileExt := strings.ToLower(path.Ext(file.Filename))
		if fileExt != ".png" && fileExt != ".jpg" && fileExt != ".gif" && fileExt != ".jpeg" {
			c.JSON(200, gin.H{
				"code": 400,
				"msg":  "上传失败!只允许png,jpg,gif,jpeg文件",
			})
			return
		}
		timeStamp := time.Now().Unix()
		fileName := fmt.Sprintf("%v%s", timeStamp, file.Filename)
		//fileDir := fmt.Sprintf("%s%s", "./image/", fileName)
		fileDir := path.Join("./image/", fileName)
		err := c.SaveUploadedFile(file, fileDir)
		if err != nil {
			fmt.Printf("%v\n", err)
		}
		c.JSON(200, gin.H{
			"code": 200,
			"msg":  "上传成功!",
			"result": gin.H{
				"path": "http://" + c.Request.Host + "/" + fileDir,
			},
		})
	}
}
