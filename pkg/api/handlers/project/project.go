package project

import (
	"fmt"
	"gorm.io/gorm/clause"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/guidewire/fern-reporter/pkg/models"
	"gorm.io/gorm"
)

type ProjectHandler struct {
	db *gorm.DB
}

// NewProjectHandler initializes ProjectHandler
func NewProjectHandler(db *gorm.DB) *ProjectHandler {
	return &ProjectHandler{db: db}
}

func (h *ProjectHandler) GetAllProjects(c *gin.Context) {
	var projects []models.ProjectDetails

	if err := h.db.Find(&projects).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching projects"})
		return
	}
	c.JSON(http.StatusOK, projects)
}

func (h *ProjectHandler) CreateProject(c *gin.Context) {
	var project models.ProjectDetails

	if err := c.ShouldBindJSON(&project); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.db.Clauses(clause.Returning{}).Create(&project).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create project"})
		return
	}

	c.JSON(http.StatusCreated, project)
}

func (h *ProjectHandler) UpdateProject(c *gin.Context) {
	id := c.Param("uuid")
	var project models.ProjectDetails

	if err := h.db.Where("uuid = ?", id).First(&project).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	if err := c.ShouldBindJSON(&project); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.db.Save(&project).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update project"})
		return
	}

	c.JSON(http.StatusOK, project)
}

func (h *ProjectHandler) DeleteProject(c *gin.Context) {
	uuid := c.Param("uuid")

	if err := h.db.Where("uuid = ?", uuid).Delete(&models.ProjectDetails{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete project"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Project ID %s deleted", uuid)})
}

func (h *ProjectHandler) GetProjectAll(c *gin.Context) {
	var projects []struct {
		ID   uint64 `json:"id"`
		Name string `json:"name"`
		UUID string `json:"uuid"`
	}
	h.db.Table("project_details").
		Order("name ASC").
		Find(&projects)

	c.JSON(http.StatusOK, gin.H{
		"projects": projects,
	})
}
