package controller

import (
	"bluebell_backend/dao/mysql"
	"bluebell_backend/dao/redis"
	"bluebell_backend/logic"
	"bluebell_backend/models"
	"fmt"
	"path"
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
	userID, err := getCurrentUserID(c)
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
	posts := redis.GetPost(order, pageNum)
	fmt.Println(len(posts))
	ResponseSuccess(c, posts)
}

func GetPostListBYUser(c *gin.Context) {
	id, _ := getCurrentUserID(c)
	postIDList := mysql.GetPostByUser(id)
	keys := make([]string, 0)
	for _, s := range postIDList {
		keys = append(keys, strconv.FormatUint(s, 10))
	}
	postList := redis.GetPostBYKeys(keys)
	ResponseSuccess(c, postList)
}

func GetOtherUserPost(c *gin.Context) {
	id := c.Param("id")
	atom, _ := strconv.Atoi(id)
	u := uint64(atom)
	postIDList := mysql.GetPostByUser(u)
	keys := make([]string, 0)
	for _, s := range postIDList {
		keys = append(keys, strconv.FormatUint(s, 10))
	}
	postList := redis.GetPostBYKeys(keys)
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
	fmt.Println(len(posts))
	ResponseSuccess(c, posts)
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

	ResponseSuccess(c, post)
}

func UploadImg(c *gin.Context) {
	f, err := c.FormFile("imgfile")
	if err != nil {
		c.JSON(200, gin.H{
			"code": 400,
			"msg":  "上传失败!",
		})
		return
	} else {

		fileExt := strings.ToLower(path.Ext(f.Filename))
		if fileExt != ".png" && fileExt != ".jpg" && fileExt != ".gif" && fileExt != ".jpeg" {
			c.JSON(200, gin.H{
				"code": 400,
				"msg":  "上传失败!只允许png,jpg,gif,jpeg文件",
			})
			return
		}
		timeStamp := time.Now().Unix()
		fileName := fmt.Sprintf("%v%s", timeStamp, f.Filename)
		//fileDir := fmt.Sprintf("%s%s", "./image/", fileName)
		fileDir := path.Join("./image/", fileName)
		err := c.SaveUploadedFile(f, fileDir)
		if err != nil {
			fmt.Printf("%v\n", err)
		}
		c.JSON(200, gin.H{
			"code": 200,
			"msg":  "上传成功!",
			"result": gin.H{
				"path": fileDir,
			},
		})
	}
}
