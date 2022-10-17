package controllers

import (
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/panhdjf/server_management_system/models"
	"github.com/rs/xid"

	"gorm.io/gorm"
)

type SeverController struct {
	DB *gorm.DB
}

func NewServerController(DB *gorm.DB) SeverController {
	return SeverController{DB}
}

func (sc *SeverController) CreateServer(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(models.User)
	var payload *models.CreateServerRequest

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	now := time.Now()
	newServer := models.Server{
		ID:          xid.New().String(),
		Name:        payload.Name,
		Status:      payload.Status,
		Ipv4:        payload.Ipv4,
		User:        currentUser.ID,
		CreatedTime: now,
		LastUpdated: now,
	}

	result := sc.DB.Create(&newServer)
	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "duplicate key") {
			ctx.JSON(http.StatusConflict, gin.H{"status": "fail", "message": "Server with that Name already exists"})
			return
		}
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": result.Error.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "data": newServer})
}

// Review Server
func (sc *SeverController) ViewServers(ctx *gin.Context) {
	var page = ctx.DefaultQuery("page", "1")
	var limit = ctx.DefaultQuery("limit", "10")

	intPage, _ := strconv.Atoi(page)
	intLimit, _ := strconv.Atoi(limit)
	offset := (intPage - 1) * intLimit
	var servers []models.Server
	serverIpv4 := ctx.Param("serverIpv4")
	results := sc.DB.Limit(intLimit).Offset(offset).Find(&servers, "Ipv4=?", serverIpv4)
	if results.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": results.Error})
		return
	}
	SortSever(servers)
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "Filter Ipv4", "results": len(servers), "data": servers})
}

func SortSever(servers []models.Server) []models.Server {
	sort.Slice(servers, func(i, j int) bool {
		return servers[i].Name > servers[j].Name
	})
	return servers
}

func (sc *SeverController) UpdateServer(ctx *gin.Context) {
	serverId := ctx.Param("serverId")
	currentUser := ctx.MustGet("currentUser").(models.User)

	var payload *models.UpdateServer
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}
	var updatedServer models.Server
	result := sc.DB.First(&updatedServer, "sever_id = ?", serverId)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "No Server with that Id exists"})
		return
	}
	now := time.Now()
	serverToUpdate := models.Server{
		ID:          updatedServer.ID,
		Name:        payload.Name,
		Status:      payload.Status,
		Ipv4:        payload.Ipv4,
		User:        currentUser.ID,
		CreatedTime: updatedServer.CreatedTime,
		LastUpdated: now,
	}

	sc.DB.Model(&updatedServer).Updates(serverToUpdate)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": updatedServer})
}

func (sc *SeverController) DeletePost(ctx *gin.Context) {
	serverId := ctx.Param("serverId")

	result := sc.DB.Delete(&models.Server{}, "sever_id = ?", serverId)

	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "No Server with that Id exists"})
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}
