package handler

import (
	"fmt"
	"net/http"

	"github.com/RomanshkVolkov/test-api/internal/core/domain"
	"github.com/RomanshkVolkov/test-api/internal/core/service"
	"github.com/gin-gonic/gin"
)

func GetDocuments(c *gin.Context) {
	server := service.GetServer(c)
	documents := server.GetDocuments()

	c.IndentedJSON(http.StatusOK, documents)
}

func GetDocumentByID(c *gin.Context) {
	server := service.GetServer(c)
	id, err := ExtractAndParseUintParam(c, "id")
	if err != nil {
		return
	}

	document := server.GetDocumentByID(id)

	c.IndentedJSON(http.StatusOK, document)
}

func UpdateDocument(c *gin.Context) {
	request, err := ValidateRequest[domain.EditableDocument](c)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, ServerError(err, RequestError))
		return
	}

	id, err := ExtractAndParseUintParam(c, "id")
	if err != nil {
		InvalidParamError(c, "id", err)
		return
	}

	if id != request.ID {
		InvalidParamError(c, "id", nil)
		return
	}

	server := service.GetServer(c)

	updatedDocument := server.UpdateDocument(request)
	// update document
	c.IndentedJSON(http.StatusOK, updatedDocument)
}

func DeleteDocument(c *gin.Context) {
	server := service.GetServer(c)
	id, err := ExtractAndParseUintParam(c, "id")
	if err != nil {
		return
	}

	response := server.DeleteDocument(id)

	c.IndentedJSON(http.StatusOK, response)
}

func DeleteFieldDocument(c *gin.Context) {
	server := service.GetServer(c)
	id, err := ExtractAndParseUintParam(c, "id")
	if err != nil {
		return
	}

	response := server.DeleteFieldDocument(id)

	c.IndentedJSON(http.StatusOK, response)
}

func CreateDocument(c *gin.Context) {
	request, err := ValidateRequest[domain.DocumentWithDetails](c)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, ServerError(err, RequestError))
		return
	}
	server := service.GetServer(c)

	createdDocument := server.CreateDocument(request)
	// create document
	c.IndentedJSON(http.StatusOK, createdDocument)
}

func GetDocumentRowRecord(c *gin.Context) {
	server := service.GetServer(c)

	fmt.Println("GetDocumentRowRecord")
	documentID, err := ExtractAndParseUintParam(c, "document")
	if err != nil {
		return
	}
	id, err := ExtractAndParseUintParam(c, "id")
	if err != nil {
		return
	}

	records := server.GetDocumentRowRecord(id, documentID)

	c.IndentedJSON(http.StatusOK, records)
}

func GetReports(c *gin.Context) {
	server := service.GetServer(c)
	reports := server.GetReports()

	c.IndentedJSON(http.StatusOK, reports)
}

func GetReportByID(c *gin.Context) {
	server := service.GetServer(c)
	id, err := ExtractAndParseUintParam(c, "id")
	if err != nil {
		return
	}

	report := server.GetReportByID(id)

	c.IndentedJSON(http.StatusOK, report)
}

func UpdateDocumentRowRecord(c *gin.Context) {
	request, err := ValidateRequest[domain.EditableDoumentRowRecord](c)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, ServerError(err, RequestError))
		return
	}
	server := service.GetServer(c)

	updatedRecords := server.UpdateDocumentRowRecord(request)
	// update records
	c.IndentedJSON(http.StatusOK, updatedRecords)
}

func DeleteDocumentRowRecord(c *gin.Context) {
	server := service.GetServer(c)
	documentID, err := ExtractAndParseUintParam(c, "document")
	if err != nil {
		return
	}
	id, err := ExtractAndParseUintParam(c, "id")
	if err != nil {
		return
	}

	response := server.DeleteDocumentRowRecord(id, documentID)

	c.IndentedJSON(http.StatusOK, response)
}

func GetTables(c *gin.Context) {
	server := service.GetServer(c)
	reports := server.GetTables()

	c.IndentedJSON(http.StatusOK, reports)
}

func GetTableByID(c *gin.Context) {
	server := service.GetServer(c)
	id, err := ExtractAndParseUintParam(c, "id")
	if err != nil {
		return
	}

	report := server.GetTableByID(id)

	c.IndentedJSON(http.StatusOK, report)
}

func UploadDocument(c *gin.Context) {
	request, err := ValidateRequest[domain.UploadDocument](c)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, ServerError(err, RequestError))
		return
	}
	server := service.GetServer(c)

	uploadedDocument := server.UploadDocument(request)
	// upload document
	c.IndentedJSON(http.StatusOK, uploadedDocument)
}
