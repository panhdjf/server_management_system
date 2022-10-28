package controllers

import (
	"encoding/json"
	"io/ioutil"
	"log"
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

type ServerController struct {
	DB *gorm.DB
}

func NewServerController(DB *gorm.DB) ServerController {
	return ServerController{DB}
}

func (sc *ServerController) CreateServer(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(models.User)
	var payload *models.CreateServerRequest

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	now := time.Now()
	newServer := models.Server{
		ID:            xid.New().String(),
		Name:          payload.Name,
		Status:        payload.Status,
		Uptime:        payload.Uptime,
		Ipv4:          payload.Ipv4,
		IdUserManager: currentUser.ID,
		CreatedTime:   now,
		LastUpdated:   now,
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

func (sc *ServerController) ViewServers(ctx *gin.Context) {
	var page = ctx.DefaultQuery("page", "0")
	var limit = ctx.DefaultQuery("limit", "10")

	intPage, _ := strconv.Atoi(page)
	intLimit, _ := strconv.Atoi(limit)
	offset := (intPage - 1) * intLimit

	var sortRequired = ctx.DefaultQuery("sortRequired", "name")
	filterRequired := ctx.DefaultQuery("filterRequired", "")
	var valueRequired = ctx.DefaultQuery("valueRequired", "")

	var servers []models.Server

	if filterRequired == "" || valueRequired == "" {
		result := sc.DB.Limit(intLimit).Offset(offset).Order(sortRequired).Find(&servers)
		if result.Error != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "No Server with that required exists"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"status": "success", "sort": sortRequired, "filter": "No filter", "results": len(servers), "data": servers})
		return
	}

	result := sc.DB.Order(sortRequired).Limit(intLimit).Offset(offset).Where(filterRequired, valueRequired).Find(&servers)

	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "No Server with that required exists"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "sort": sortRequired, "filter": filterRequired, "required filter": valueRequired, "results": len(servers), "data": servers})
}

func (sc *ServerController) UpdateServer(ctx *gin.Context) {
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
		ID:            updatedServer.ID,
		Name:          payload.Name,
		Status:        payload.Status,
		Uptime:        payload.Uptime,
		Ipv4:          payload.Ipv4,
		IdUserManager: currentUser.ID,
		CreatedTime:   updatedServer.CreatedTime,
		LastUpdated:   now,
	}

	sc.DB.Model(&updatedServer).Updates(serverToUpdate)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": updatedServer})
}

func (sc *ServerController) DeleteServer(ctx *gin.Context) {
	serverId := ctx.Param("serverId")
	var server models.Server
	result := sc.DB.Delete(&server, "id = ?", serverId)

	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "No Server with that Id exists"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "deletedsuccess"})
}

func (sc *ServerController) DeleteAllServers(ctx *gin.Context) {
	var servers []models.Server
	sc.DB.Find(&servers)

	results := sc.DB.Delete(&servers)
	if results.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": results.Error.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"results": "All data have been deleted successfully"})
}

func (sc *ServerController) ExportExcel(ctx *gin.Context) {
	f := excelize.NewFile()
	// Create a new sheet.
	// index := f.NewSheet("Sheet1")
	// Set value of a cell.

	f.SetCellValue("Sheet1", "A1", "ID")
	f.SetCellValue("Sheet1", "B1", "Name")
	f.SetCellValue("Sheet1", "C1", "Status")
	f.SetCellValue("Sheet1", "D1", "Uptime")
	f.SetCellValue("Sheet1", "E1", "Ipv4")
	f.SetCellValue("Sheet1", "F1", "IdUserManager")
	f.SetCellValue("Sheet1", "G1", "CreatedTime")
	f.SetCellValue("Sheet1", "H1", "LastUpdated")
	var Servers []models.Server
	//get servers from DB
	//ex: sort with name
	var sortRequired = ctx.DefaultQuery("sortRequired", "name")

	sc.DB.Order(sortRequired).Find(&Servers)

	for i, server := range Servers {
		f.SetCellValue("Sheet1", "A"+strconv.Itoa(i+2), server.ID)
		f.SetCellValue("Sheet1", "B"+strconv.Itoa(i+2), server.Name)
		f.SetCellValue("Sheet1", "C"+strconv.Itoa(i+2), server.Status)
		f.SetCellValue("Sheet1", "D"+strconv.Itoa(i+2), server.Uptime)
		f.SetCellValue("Sheet1", "E"+strconv.Itoa(i+2), server.Ipv4)
		f.SetCellValue("Sheet1", "F"+strconv.Itoa(i+2), server.IdUserManager)
		f.SetCellValue("Sheet1", "G"+strconv.Itoa(i+2), server.CreatedTime)
		f.SetCellValue("Sheet1", "H"+strconv.Itoa(i+2), server.LastUpdated)
		// Set active sheet of the workbook.
	}
	// f.SetActiveSheet(index)
	// Save xlsx file by the given path.
	if err := f.SaveAs("ExportServer.xlsx"); err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "Failed to export Database to the excel", "error": err.Error()})
		return
	}
	// ctx.FileAttachment("./C:/Users/sv_anhntp/PA/server_management_system", "ExportServer.xlsx")
	// ctx.Header("Conte0nt-Disposition" `attachment; filename="gopher.png"`)
	// ctx.File("./Users/sv_anhntp/PA/server_management_system/ExportServer.xlsx")

	ctx.JSON(http.StatusCreated, gin.H{"status": "success"})
}

func (sc *ServerController) ImportExcel(ctx *gin.Context) {

	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "Failed to import Database to the excel", "error": err.Error()})
		return
	}

	f, err := excelize.OpenFile(file.Filename)
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
	var servers []models.Server
	sc.DB.Find(&servers)

	serversImport := make([]models.Server, 0)

	serversAccept := make([]models.ImportExcel, 0)
	serversFail := make([]models.ImportExcel, 0)

	if len(servers) != 0 {
		for _, server := range servers {
			for _, row := range rows {
				if len(row) != 0 {
					if server.ID == row[0] || server.Name == row[1] {
						newServerFail := models.ImportExcel{
							ID:   row[0],
							Name: row[1],
						}
						serversFail = append(serversFail, newServerFail)
						continue
					}
					user, _ := strconv.Atoi(row[4])
					uptime, _ := strconv.ParseFloat(row[3], 8)
					newServer := models.Server{
						ID:            row[0],
						Name:          row[1],
						Status:        row[2],
						Uptime:        uptime,
						Ipv4:          row[4],
						IdUserManager: user,
						CreatedTime:   now,
						LastUpdated:   now,
					}
					serversImport = append(serversImport, newServer)

					newServerAccept := models.ImportExcel{
						ID:   row[0],
						Name: row[1],
					}
					serversAccept = append(serversAccept, newServerAccept)
				}
			}
		}
	} else {
		for _, row := range rows {
			if len(row) != 0 {
				user, _ := strconv.Atoi(row[4])
				uptime, _ := strconv.ParseFloat(row[3], 8)
				newServer := models.Server{
					ID:            row[0],
					Name:          row[1],
					Status:        row[2],
					Uptime:        uptime,
					Ipv4:          row[4],
					IdUserManager: user,
					CreatedTime:   now,
					LastUpdated:   now,
				}
				serversImport = append(serversImport, newServer)

				newServerAccept := models.ImportExcel{
					ID:   row[0],
					Name: row[1],
				}
				serversAccept = append(serversAccept, newServerAccept)
			}
		}
	}
	results := sc.DB.Create(&serversImport)

	if results.Error != nil {
		ctx.JSON(http.StatusOK, gin.H{"status": "fail", "message": results.Error.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"status": gin.H{"ImportEccept": gin.H{"CountAccept": len(serversAccept), "data": serversAccept}, "ImportFail": gin.H{"CountFail": len(serversFail), "data": serversFail}}})
}

func (sc ServerController) CheckStatusServer() (int, int, int, float64) {

	var servers []models.Server
	sc.DB.Find(&servers)
	totalServer := len(servers)
	if totalServer == 0 {
		log.Fatal("No server exists")
	}
	countServerOn := 0
	countServerOff := 0
	totalUptime := 0.0
	for _, server := range servers {
		totalUptime += server.Uptime
		if server.Status == "offline" {
			countServerOff++
			continue
		}
		countServerOn++
	}

	var avgUptime float64
	avgUptime = totalUptime / float64(totalServer)
	return totalServer, countServerOn, countServerOff, avgUptime
}

func (sc ServerController) UpdateStatusServer(ctx *gin.Context) {
	var servers []models.Server
	sc.DB.Find(&servers)
	totalServer := len(servers)
	if totalServer == 0 {
		log.Fatal("No server exists")
	}
	now := time.Now()
	for _, server := range servers {
		var status string
		var uptime float64
		url := strings.Join([]string{"http://", server.ID, ":8000/status"}, "")
		response, err := http.Get(url)
		// response, err := http.Get("http://192.168.2.0:8000/status")
		if err != nil {
			status = "offline"
			continue
		}
		status = "online"
		responseData, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}
		var responseServer models.ServerStatus
		err1 := json.Unmarshal(responseData, &responseServer)
		if err1 != nil {
			log.Fatal(err1)
		}
		uptime, _ = strconv.ParseFloat(responseServer.UpdateTime, 8)
		// fmt.Println(responseServer.UpdateTime)
		// fmt.Println(responseServer.Status)

		serverToUpdate := models.Server{
			ID:            server.ID,
			Name:          server.Name,
			Status:        status,
			Uptime:        uptime,
			Ipv4:          server.Ipv4,
			IdUserManager: server.IdUserManager,
			CreatedTime:   server.CreatedTime,
			LastUpdated:   now,
		}

		sc.DB.Model(&server).Updates(serverToUpdate)
	}
}
