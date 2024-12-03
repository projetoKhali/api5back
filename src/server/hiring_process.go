package server

import (
	"net/http"

	"api5back/ent"
	"api5back/src/model"
	"api5back/src/service"

	"github.com/gin-gonic/gin"
)

func HiringProcessDashboard(
	engine *gin.Engine,
	dbClient *ent.Client,
	dwClient *ent.Client,
) {
	v1 := engine.Group("/api/v1")
	{
		hiringProcess := v1.Group("/hiring-process")
		{
			hiringProcess.POST("/dashboard", Dashboard(dwClient))
			hiringProcess.POST("/table", VacancyTable(dwClient))
		}

		suggestions := v1.Group("/suggestions")
		{
			suggestions.POST("/recruiter", UserList(dwClient))
			suggestions.POST("/process", HiringProcessList((dwClient)))
			suggestions.POST("/vacancy", VacancyList(dwClient))
			suggestions.GET("/department", ListDepartments(dbClient))
		}
		authentication := v1.Group("/authentication")
		{
			authentication.GET("/users", ListUsers(dbClient))
			authentication.POST("/login", LoginUser(dbClient))
			authentication.POST("/create", CreateUser(dbClient))
		}
		accessGroup := v1.Group("/access-group")
		{
			accessGroup.GET("", ListAccessGroup(dbClient))
			accessGroup.POST("", CreateAccessGroup(dbClient))
		}
	}
}

// Dashboard godoc
// @Summary dashboard
// @Description show dashboard
// @Tags hiring-process
// @Accept json
// @Param body body model.FactHiringProcessFilter true "Metrics filter"
// @Produce json
// @Success 200 {string} Dashboard
// @Router /hiring-process/dashboard [post]
func Dashboard(
	dwClient *ent.Client,
) func(c *gin.Context) {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")

		// TODO: change to pointer
		var dashboardMetricsFilter model.FactHiringProcessFilter
		if err := c.ShouldBindJSON(&dashboardMetricsFilter); err != nil {
			c.JSON(http.StatusBadRequest, DisplayError(err))
			return
		}

		metricsData, err := service.GetMetrics(
			c, dwClient, dashboardMetricsFilter,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, DisplayError(err))
			return
		}

		c.JSON(http.StatusOK, metricsData)
	}
}

// UserList godoc
// @Summary List users
// @Description Return a list of users with id and name
// @Tags suggestions
// @Accept json
// @Param body body model.PageRequest true "Page request"
// @Produce json
// @Success 200 {array} model.Page[model.Suggestion]
// @Router /suggestions/recruiter/ [post]
func UserList(dwClient *ent.Client) func(c *gin.Context) {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")

		var pageRequest *model.SuggestionsPageRequest

		if err := c.ShouldBindJSON(&pageRequest); err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}

		users, err := service.GetUserSuggestions(
			c, dwClient,
			pageRequest,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, DisplayError(err))
			return
		}

		c.JSON(http.StatusOK, users)
	}
}

// HiringProcessList godoc
// @Summary List hiring processes
// @Description Return a list of hiring processes with id and title
// @Tags suggestions
// @Accept json
// @Param body body model.SuggestionsFilter true "Filter"
// @Produce json
// @Success 200 {array} model.Page[model.Suggestion]
// @Router /suggestions/process [post]
func HiringProcessList(
	dbClient *ent.Client,
) func(c *gin.Context) {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")

		var pageRequest *model.SuggestionsFilter

		if err := c.ShouldBindJSON(&pageRequest); err != nil {
			c.JSON(http.StatusBadRequest, DisplayError(err))
			return
		}

		processes, err := service.GetProcessSuggestions(
			c, dbClient,
			pageRequest,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, DisplayError(err))
			return
		}

		c.JSON(http.StatusOK, processes)
	}
}

// HiringProcessList godoc
// @Summary List hiring processes
// @Description Return a list of hiring processes with id and title
// @Tags suggestions
// @Accept json
// @Param body body model.SuggestionsFilter true "Filter"
// @Produce json
// @Success 200 {array} model.Page[model.Suggestion]
// @Router /suggestions/vacancy [post]
func VacancyList(
	dwClient *ent.Client,
) func(c *gin.Context) {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")

		var pageRequest *model.SuggestionsFilter

		if err := c.ShouldBindJSON(&pageRequest); err != nil {
			c.JSON(http.StatusBadRequest, DisplayError(err))
			return
		}

		vacancies, err := service.GetVacancySuggestions(
			c, dwClient,
			pageRequest,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, DisplayError(err))
			return
		}

		c.JSON(http.StatusOK, vacancies)
	}
}

// HiringProcessList godoc
// @Summary List hiring processes
// @Description Return a list of vacancies with summarized information
// @Tags hiring-process
// @Accept json
// @Param body body model.FactHiringProcessFilter true "Metrics filter"
// @Produce json
// @Success 200 {array} model.Page[model.DashboardTableRow]
// @Router /hiring-process/table [post]
func VacancyTable(
	dwClient *ent.Client,
) func(c *gin.Context) {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")

		var filter model.FactHiringProcessFilter

		if err := c.ShouldBindJSON(&filter); err != nil {
			c.JSON(http.StatusBadRequest, DisplayError(err))
			return
		}

		vacancies, err := service.GetVacancyTable(
			c, dwClient,
			filter,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, DisplayError(err))
			return
		}

		c.JSON(http.StatusOK, vacancies)
	}
}

// ListDepartments godoc
// @Summary List departments
// @Description Return a list of departments with id and title
// @Tags departments
// @Produce json
// @Success 200 {array} model.Suggestion
// @Router /suggestions/departments [get]
func ListDepartments(
	client *ent.Client,
) func(c *gin.Context) {
	return func(c *gin.Context) {
		departments, err := service.ListDepartments(c, client)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to list departments",
				"details": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, departments)
	}
}

// CreateAccessGroup godoc
// @Summary Create a new group access
// @Description Create a new group access with name and related departments
// @Tags access_group
// @Accept json
// @Produce json
// @Param body body model.CreateAccessGroupRequest true "Group Access Info"
// @Success 201 {object} ent.AccessGroup
// @Router /access-group [post]
func CreateAccessGroup(client *ent.Client) func(c *gin.Context) {
	return func(c *gin.Context) {
		var request model.CreateAccessGroupRequest

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		group, err := service.CreateAccessGroup(c, client, request)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, group)
	}
}

// ListAccessGroup godoc
// @Summary List group access with departments
// @Description Return a list of group access with id, name, and departments
// @Tags access_group
// @Produce json
// @Success 200 {array} model.AccessGroup
// @Router /access-group [get]
func ListAccessGroup(client *ent.Client) func(c *gin.Context) {
	return func(c *gin.Context) {
		groups, err := service.GetAccessGroupWithDepartments(c, client)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, groups)
	}
}

// ListUsers godoc
// @Summary List users
// @Description Return a list of users with name, email, and group
// @Tags authentication
// @Produce json
// @Success 200 {array} service.UserResponse
// @Router /authentication/users [get]
func ListUsers(client *ent.Client) func(c *gin.Context) {
	return func(c *gin.Context) {
		users, err := client.Authentication.Query().
			WithAccessGroup().
			All(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		var response []service.UserResponse
		for _, user := range users {
			response = append(response, service.UserResponse{
				Name:  user.Name,
				Email: user.Email,
				Group: user.Edges.AccessGroup.Name,
			})
		}

		c.JSON(http.StatusOK, response)
	}
}

// LoginUser godoc
// @Summary User login
// @Description Authenticate a user with email and password
// @Tags authentication
// @Accept json
// @Produce json
// @Param body body map[string]string true "User credentials: {email, password}"
// @Router /authentication/login [post]
func LoginUser(client *ent.Client) func(c *gin.Context) {
	return func(c *gin.Context) {
		var creds map[string]string
		if err := c.ShouldBindJSON(&creds); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}
		email, emailOk := creds["email"]
		password, passOk := creds["password"]
		if !emailOk || !passOk {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing email or password"})
			return
		}

		loginResponse, err := service.Login(c.Request.Context(), client, service.LoginRequest{
			Email:    email,
			Password: password,
		})
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "login successful",
			"user":    loginResponse,
		})
	}
}

// CreateUser godoc
// @Summary Create a new user
// @Description Create a new user with name, email, password, and group ID
// @Tags authentication
// @Accept json
// @Produce json
// @Param body body service.CreateUserRequest true "User info: {name, email, password, groupID}"
// @Success 201 {object} ent.Authentication
// @Router /authentication/create [post]
func CreateUser(client *ent.Client) func(c *gin.Context) {
	return func(c *gin.Context) {
		var request service.CreateUserRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		user, err := service.CreateUser(c, client, request)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, user)
	}
}
