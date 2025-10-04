package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

func SetupApiRouter(router *mux.Router) {
	apiRouter := router.
		PathPrefix("/api").
		Subrouter()
	inventoryRouter := apiRouter.
		PathPrefix("/inventory").
		Subrouter()
	projectBaseRouter := apiRouter.
		PathPrefix("/project").
		Subrouter()
	instructionRouter := apiRouter.
		PathPrefix("/instruction").
		Subrouter()

	inventoryBoxRouter := inventoryRouter.
		PathPrefix("/box").
		Subrouter()
	inventoryBoxRouter.
		Methods(http.MethodGet).
		HandlerFunc(getInventoryBoxes)
	inventoryBoxRouter.
		Methods(http.MethodPost).
		HandlerFunc(createInventoryBox)
	inventoryBoxRouter.
		Methods(http.MethodGet).
		Path("/{boxId}").
		HandlerFunc(getInventoryBox)

	inventoryItemRouter := inventoryBoxRouter.
		PathPrefix("/{boxId}/item").
		Subrouter()
	inventoryItemRouter.
		Methods(http.MethodGet).
		HandlerFunc(getInventoryItems)
	inventoryItemRouter.
		Methods(http.MethodPost).
		HandlerFunc(createInventoryItem)
	inventoryItemRouter.
		Methods(http.MethodGet).
		Path("/{itemId}").
		HandlerFunc(getInventoryItem)
	inventoryItemRouter.
		Methods(http.MethodPut).
		Path("/{itemId}").
		HandlerFunc(updateInventoryItem)
	inventoryItemRouter.
		Methods(http.MethodDelete).
		Path("/{itemId}").
		HandlerFunc(deleteInventoryItem)

	projectCategoryRouter := projectBaseRouter.
		PathPrefix("/category").
		Subrouter()
	projectCategoryRouter.
		Methods(http.MethodGet).
		HandlerFunc(getProjectCategories)
	projectCategoryRouter.
		Methods(http.MethodPost).
		HandlerFunc(createProjectCategory)
	projectCategoryRouter.
		Methods(http.MethodGet).
		Path("/{categoryId}").
		HandlerFunc(getProjectCategory)

	projectRouter := projectCategoryRouter.
		PathPrefix("/{categoryId}/project").
		Subrouter()
	projectRouter.
		Methods(http.MethodGet).
		HandlerFunc(getProjects)
	projectRouter.
		Methods(http.MethodPost).
		HandlerFunc(createProject)
	projectRouter.
		Methods(http.MethodGet).
		Path("/{projectId}").
		HandlerFunc(getProject)
	projectRouter.
		Methods(http.MethodPut).
		Path("/{projectId}").
		HandlerFunc(updateProject)
	projectRouter.
		Methods(http.MethodDelete).
		Path("/{projectId}").
		HandlerFunc(deleteProject)
	projectRouter.
		Methods(http.MethodPut).
		Path("/{projectId}/archive").
		HandlerFunc(archiveProject)
	projectRouter.
		Methods(http.MethodDelete).
		Path("/{projectId}/archive").
		HandlerFunc(unarchiveProject)

	instructionRouter.
		Methods(http.MethodGet).
		HandlerFunc(getInstructions)
	instructionRouter.
		Methods(http.MethodGet).
		Path("/{instructionId}").
		HandlerFunc(getInstruction)
	instructionRouter.
		Methods(http.MethodDelete).
		Path("/{instructionId}").
		HandlerFunc(deleteInstruction)
	instructionRouter.
		Methods(http.MethodPut).
		Path("/{instructionId}").
		HandlerFunc(updateInstruction)
	instructionRouter.
		Methods(http.MethodPut).
		Path("/{instructionId}/steps").
		HandlerFunc(replaceSteps)

	instructionStepRouter := instructionRouter.
		Path("/{instructionId}/step").
		Subrouter()
	instructionStepRouter.
		Methods(http.MethodGet).
		HandlerFunc(getInstructionSteps)
	instructionStepRouter.
		Methods(http.MethodPost).
		HandlerFunc(createInstructionStep)
	instructionRouter.
		Methods(http.MethodGet).
		Path("/{stepId}").
		HandlerFunc(getInstructionStep)
	instructionRouter.
		Methods(http.MethodPut).
		Path("/{stepId}").
		HandlerFunc(updateInstructionStep)
	instructionRouter.
		Methods(http.MethodDelete).
		Path("/{stepId}").
		HandlerFunc(deleteInstructionStep)
	instructionRouter.
		Methods(http.MethodPut).
		Path("/{stepId}/done").
		HandlerFunc(markInstructionStepAsDone)
	instructionRouter.
		Methods(http.MethodDelete).
		Path("/{stepId}/done").
		HandlerFunc(markInstructionStepAsTodo)

	apiRouter.
		Methods(http.MethodGet).
		Path("/healthz").
		HandlerFunc(getHealth)

	apiRouter.Use(authentication(), contentTypeJson)
}
