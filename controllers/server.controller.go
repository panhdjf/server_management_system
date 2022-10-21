package controllers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/panhdjf/server_management_system/models"
	"github.com/rs/xid"
	"github.com/xuri/excelize/v2"

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

func (sc *SeverController) ViewServers(ctx *gin.Context) {
	var page = ctx.DefaultQuery("page", "1")
	var limit = ctx.DefaultQuery("limit", "10")

	intPage, _ := strconv.Atoi(page)
	intLimit, _ := strconv.Atoi(limit)
	offset := (intPage - 1) * intLimit

	var requiredFilter = ctx.DefaultQuery("requiredFilter", "")
	var requiredSort = ctx.DefaultQuery("requiredSort", "name")

	var servers []models.Server
	result := sc.DB.Limit(intLimit).Offset(offset).Order(requiredSort).Find(&servers, "status = ?", requiredFilter)

	// switch required {
	// case "status":
	// 	{
	// 		result := sc.DB.Limit(intLimit).Offset(offset).Find(&servers, "status = ?", valueRequired)
	// 	}
	// case "name":
	// 	{
	// 		result := sc.DB.Limit(intLimit).Offset(offset).Find(&servers, "name = ?", valueRequired)
	// 	}
	// }

	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "No Server with that required exists"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "results": len(servers), "data": servers})
}

// func SortSever(servers []models.Server, required string) []models.Server {
// 	switch required {
// 	case "serverName":
// 		sort.Slice(servers, func(i, j int) bool {
// 			return servers[i].Name < servers[j].Name
// 		})
// 	case "serverIpv4":
// 		sort.Slice(servers, func(i, j int) bool {
// 			return servers[i].Ipv4 < servers[j].Ipv4
// 		})

// 	case "serverUser":
// 		sort.Slice(servers, func(i, j int) bool {
// 			return servers[i].User < servers[j].User
// 		})

// 	case "serverStatus":
// 		sort.Slice(servers, func(i, j int) bool {
// 			return servers[i].Status < servers[j].Status
// 		})

// 	// case "sort by created time":
// 	// 	sort.Slice(servers, func(i, j int) bool {
// 	// 		return servers[i].CreatedTime < servers[j].CreatedTime
// 	// 	})

// 	// case "sort by last updated time":
// 	// 	sort.Slice(servers, func(i, j int) bool {
// 	// 		return servers[i].LastUpdated < servers[j].LastUpdated
// 	// 	})

// 	default:
// 		log.Fatal("Error")
// 	}

// 	return servers
// }

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

func (sc *SeverController) DeleteAllServers(ctx *gin.Context) {
	var servers []models.Server
	sc.DB.Offset(0).Find(&servers)

	results := sc.DB.Delete(&servers)
	if results.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": results.Error.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"results": "All data have been deleted successfully"})
}

func (sc *SeverController) ExportExcel(ctx *gin.Context) {
	f := excelize.NewFile()
	// Create a new sheet.
	// index := f.NewSheet("Sheet1")

	// Set value of a cell.
	var Servers []models.Server
	//get servers from DB
	sc.DB.Offset(0).Find(&Servers)
	for i, server := range Servers {
		f.SetCellValue("Sheet1", "A"+strconv.Itoa(i+2), server.ID)
		f.SetCellValue("Sheet1", "B"+strconv.Itoa(i+2), server.Name)
		f.SetCellValue("Sheet1", "C"+strconv.Itoa(i+2), server.Status)
		f.SetCellValue("Sheet1", "D"+strconv.Itoa(i+2), server.Ipv4)
		f.SetCellValue("Sheet1", "E"+strconv.Itoa(i+2), server.CreatedTime)
		f.SetCellValue("Sheet1", "F"+strconv.Itoa(i+2), server.LastUpdated)
		// Set active sheet of the workbook.
	}
	// f.SetActiveSheet(index)
	// Save xlsx file by the given path.
	if err := f.SaveAs("Server.xlsx"); err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "Failed to export Database to the excel", "error": err})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"status": "success"})
}

func (sc *SeverController) ImportExcel(ctx *gin.Context) {
	f, err := excelize.OpenFile("Server.xlsx")
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "Failed to import Database to the excel", "error": err})
		return
	}

	rows, err := f.GetRows("Sheet1")
	ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": err})
	return
	var servers []models.Server
	sc.DB.Offset(0).Find(&servers)
	// countAccept := 0
	// countFail := 0
	now := time.Now()
	var serversAccept []models.Server
	var serversFail []models.Server
	for _, server := range servers {
		for _, row := range rows {
			if len(row) != 0 {
				newServer := models.Server{
					ID:          row[0],
					Name:        row[1],
					Status:      row[2],
					Ipv4:        row[3],
					CreatedTime: now,
					LastUpdated: now,
				}
				if server.ID == newServer.ID || server.Name == newServer.Name {
					// countFail++
					serversFail = append(serversFail, newServer)
					continue
				}
				// countAccept ++
				serversAccept = append(serversAccept, newServer)
			}
		}
	}
	results1 := sc.DB.Create(&serversAccept)
	results2 := sc.DB.Create(&serversFail)

	if results1.Error != nil || results2.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": results1.Error.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"status": gin.H{"ImportEccept": gin.H{"CountAccept": len(serversAccept), "data": serversAccept}, "ImportFail": gin.H{"CountFail": len(serversFail), "data": serversFail}}})
}
