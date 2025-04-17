package kanboard

import (
	"server/internal/app/kanboard/handlers"
	"server/internal/constant"
	"server/internal/middlewares"

	"github.com/gin-gonic/gin"
)

func InitApi(router *gin.RouterGroup) {
	kanboard := router.Group("/kanboard")

	user := kanboard.Group("/:id")
	user.Use(middlewares.Auth(constant.KANBOARD_TOKEN))
	{
		userHandler := handlers.NewUserHandler()

		kanboard.POST("/login", userHandler.Login)
		kanboard.POST("/register", userHandler.Register)

		user.GET("/getUser", userHandler.GetUserById)
		user.PUT("/update", userHandler.UpdateInfo)
		user.POST("/password", userHandler.UpdatePassword)
		user.GET("/statistics", userHandler.GetStatistics)
		user.GET("/calendar", userHandler.GetCalendar)
	}

	public := kanboard.Group("/")
	publicHandler := handlers.NewPublicHandler()
	{
		public.GET("/captcha", publicHandler.GetCaptcha)
		user.POST("/upload", publicHandler.Upload)
	}

	wsHandler := handlers.NewWSHandler()
	{
		user.GET("/ws", wsHandler.WebSocket)
	}

	messageHandler := handlers.NewMessageHandler()
	{
		user.GET("/unreadMsgs", messageHandler.GetUnReadMsgs)
		user.GET("/readedMsgs", messageHandler.GetReadedMsgs)
		user.POST("/markReadMsg", messageHandler.MarkReadMsg)
		user.POST("/getMsgsByProjectId", messageHandler.GetMsgsByProjectId)
	}

	projectHandler := handlers.NewProjectHandler()
	{

		user.GET("/project", projectHandler.GetProjectById)
		user.POST("/addProjectMember", projectHandler.AddProjectMember)
		user.POST("/removeProjectMember", projectHandler.RemoveProjectMember)
		user.POST("/getProjectMembers", projectHandler.GetProjectMembers)
		user.POST("/getMembers", projectHandler.GetMembers)
	}

	taskHandler := handlers.NewTaskHandler()
	{
		user.GET("/tasks", taskHandler.GetTaskListByProjectId)
		user.GET("/tasksByProjectId", taskHandler.GetTaskByProjectId)
		user.GET("/tasksByUserId", taskHandler.GetTaskListByUserId)
		user.GET("/userTasks", taskHandler.GetTaskByUserId)
		user.POST("/createTask", taskHandler.CreateTask)
		user.POST("/updateTask", taskHandler.UpdateTask)
		user.POST("/updateTaskStatus", taskHandler.UpdateTaskStatus)
		user.DELETE("/deleteTask", taskHandler.DeleteTask)
		user.GET("/getTask", taskHandler.GetTask)
		user.POST("/addTaskAssignee", taskHandler.AddTaskAssignee)
		user.POST("/removeTaskAssignee", taskHandler.RemoveTaskAssignee)
		user.POST("/searchTask", taskHandler.SearchTask)
	}
}
