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

func (sc *SeverController) SortServers(ctx *gin.Context) {
	var offset = ctx.DefaultQuery("offset", "0")
	var limit = ctx.DefaultQuery("limit", "10")

	intOffset, _ := strconv.Atoi(offset)
	intLimit, _ := strconv.Atoi(limit)

	var sortRequired = ctx.DefaultQuery("sortRequired", "name")
	var servers []models.Server

	// offset: bo qua offset servers dau
	//limit: lay limit servers
	//order: sap sap theo Vd: tang dan: "name", giam dan: "name decs"

	result := sc.DB.Limit(intLimit).Offset(intOffset).Order(sortRequired).Find(&servers)

	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "No Server with that required exists"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "sort": sortRequired, "results": len(servers), "data": servers})
}

func (sc *SeverController) FilterAndSortServers(ctx *gin.Context) {
	var offset = ctx.DefaultQuery("offset", "0")
	var limit = ctx.DefaultQuery("limit", "10")

	intOffset, _ := strconv.Atoi(offset)
	intLimit, _ := strconv.Atoi(limit)

	var filterRequired = ctx.DefaultQuery("filterRequired", "status")
	var valueRequired = ctx.DefaultQuery("valueRequired", "online")

	var sortRequired = ctx.DefaultQuery("sortRequired", "name")

	var servers []models.Server

	//Ex: Filter with status = online
	result := sc.DB.Order(sortRequired).Limit(intLimit).Offset(intOffset).Where(filterRequired, valueRequired).Find(&servers)

	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "No Server with that required exists"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "sort": sortRequired, "results": len(servers), "filter": filterRequired, "value": valueRequired, "data": servers})
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
	result := sc.DB.First(&updatedServer, "id = ?", serverId)
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
	var server models.Server
	result := sc.DB.Delete(&server, "id = ?", serverId)

	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "No Server with that Id exists"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "deletedsuccess"})
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

	f.SetCellValue("Sheet1", "A1", "ID")
	f.SetCellValue("Sheet1", "B1", "Name")
	f.SetCellValue("Sheet1", "C1", "Status")
	f.SetCellValue("Sheet1", "D1", "Ipv4")
	f.SetCellValue("Sheet1", "E1", "CreatedTime")
	f.SetCellValue("Sheet1", "F1", "LastUpdated")
	var Servers []models.Server
	//get servers from DB
	sc.DB.Offset(0).Find(&Servers)
	for i, server := range Servers {
		f.SetCellValue("Sheet1", "A"+strconv.Itoa(i+2), server.ID)
		f.SetCellValue("Sheet1", "B"+strconv.Itoa(i+2), server.Name)
		f.SetCellValue("Sheet1", "C"+strconv.Itoa(i+2), server.Status)
		f.SetCellValue("Sheet1", "D"+strconv.Itoa(i+2), server.Ipv4)
		f.SetCellValue("Sheet1", "E"+strconv.Itoa(i+2), server.User)
		f.SetCellValue("Sheet1", "F"+strconv.Itoa(i+2), server.CreatedTime)
		f.SetCellValue("Sheet1", "G"+strconv.Itoa(i+2), server.LastUpdated)
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
	var servers []models.Server
	sc.DB.Offset(0).Find(&servers)

	f, err := excelize.OpenFile("Server.xlsx")
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "Failed to import Database to the excel", "error": err})
		return
	}

	rows, err := f.GetRows("Sheet1")
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": err.Error()})
		return
	}

	now := time.Now()
	serversAccept := make([]models.Server, 0)
	var countFail int
	if len(servers) != 0 {
		for _, server := range servers {
			for _, row := range rows {
				if len(row) != 0 {
					if server.ID == row[0] || server.Name == row[1] {
						countFail += 1
						continue
					}
					user, _ := strconv.Atoi(row[4])
					newServer := models.Server{
						ID:          row[0],
						Name:        row[1],
						Status:      row[2],
						Ipv4:        row[3],
						User:        user,
						CreatedTime: now,
						LastUpdated: now,
					}
					serversAccept = append(serversAccept, newServer)
				}
			}
		}
	} else {
		for _, row := range rows {
			if len(row) != 0 {
				user, _ := strconv.Atoi(row[4])
				newServer := models.Server{
					ID:          row[0],
					Name:        row[1],
					Status:      row[2],
					Ipv4:        row[3],
					User:        user,
					CreatedTime: now,
					LastUpdated: now,
				}
				serversAccept = append(serversAccept, newServer)
			}
		}
	}
	results := sc.DB.Create(&serversAccept)

	if results.Error != nil {
		ctx.JSON(http.StatusOK, gin.H{"status": "fail", "message": results.Error.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"status": gin.H{"ImportEccept": gin.H{"CountAccept": len(serversAccept), "data": serversAccept}, "ImportFail": gin.H{"CountFail": countFail}}})
}
