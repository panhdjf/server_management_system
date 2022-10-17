package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/panhdjf/server_management_system/models"
	"github.com/xuri/excelize/v2"

	"gorm.io/gorm"
)

type ExcelController struct {
	DB *gorm.DB
}

func (ec *ExcelController) ExportExcel(ctx *gin.Context) {
	f := excelize.NewFile()
	// Create a new sheet.
	// index := f.NewSheet("Sheet1")

	// Set value of a cell.
	var Servers []models.Server
	//get servers from DB
	ec.DB.Offset(0).Find(&Servers)
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

func (ec *ExcelController) ImportExcel(ctx *gin.Context) {
	f, err := excelize.OpenFile("Server.xlsx")
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "Failed to import Database to the excel", "error": err})
		return
	}

	rows, err := f.GetRows("Sheet1")
	ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": err})
	return
	var servers []models.Server
	ec.DB.Offset(0).Find(&servers)
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
	results1 := ec.DB.Create(&serversAccept)
	results2 := ec.DB.Create(&serversFail)

	if results1.Error != nil || results2.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": results1.Error.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"status": gin.H{"ImportEccept": gin.H{"CountAccept": len(serversAccept), "data": serversAccept}, "ImportFail": gin.H{"CountFail": len(serversFail), "data": serversFail}}})
}
