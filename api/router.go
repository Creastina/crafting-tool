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

	inventoryRouter.
		Path("/search").
		Methods(http.MethodGet).
		HandlerFunc(searchInventory)

	inventoryBoxRouter := inventoryRouter.
		PathPrefix("/box").
		Subrouter()
	inventoryBoxRouter.
		Methods(http.MethodGet).
		Path("").
		HandlerFunc(getInventoryBoxes)
	inventoryBoxRouter.
		Methods(http.MethodPost).
		Path("").
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
		Path("").
		HandlerFunc(getInventoryItems)
	inventoryItemRouter.
		Methods(http.MethodPost).
		Path("").
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
	inventoryItemRouter.
		Methods(http.MethodPut).
		Path("/{itemId}/stock").
		HandlerFunc(increaseInventoryItemStock)
	inventoryItemRouter.
		Methods(http.MethodDelete).
		Path("/{itemId}/stock").
		HandlerFunc(decreaseInventoryItemStock)

	projectCategoryRouter := projectBaseRouter.
		PathPrefix("/category").
		Subrouter()
	projectCategoryRouter.
		Methods(http.MethodGet).
		Path("").
		HandlerFunc(getProjectCategories)
	projectCategoryRouter.
		Methods(http.MethodPost).
		Path("").
		HandlerFunc(createProjectCategory)
	projectCategoryRouter.
		Methods(http.MethodGet).
		Path("/{categoryId}").
		HandlerFunc(getProjectCategory)
	projectCategoryRouter.
		Methods(http.MethodPut).
		Path("/{categoryId}/archive").
		HandlerFunc(archiveCategory)

	projectRouter := projectCategoryRouter.
		PathPrefix("/{categoryId}/project").
		Subrouter()
	projectRouter.
		Methods(http.MethodGet).
		Path("").
		HandlerFunc(getProjects)
	projectRouter.
		Methods(http.MethodPost).
		Path("").
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
		Path("").
		HandlerFunc(getInstructions)
	instructionRouter.
		Methods(http.MethodPost).
		Path("").
		HandlerFunc(createInstruction)
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
		PathPrefix("/{instructionId}/step").
		Subrouter()
	instructionStepRouter.
		Methods(http.MethodGet).
		Path("").
		HandlerFunc(getInstructionSteps)
	instructionStepRouter.
		Methods(http.MethodPost).
		Path("").
		HandlerFunc(createInstructionStep)
	instructionStepRouter.
		Methods(http.MethodGet).
		Path("/{stepId}").
		HandlerFunc(getInstructionStep)
	instructionStepRouter.
		Methods(http.MethodPut).
		Path("/{stepId}").
		HandlerFunc(updateInstructionStep)
	instructionStepRouter.
		Methods(http.MethodDelete).
		Path("/{stepId}").
		HandlerFunc(deleteInstructionStep)
	instructionStepRouter.
		Methods(http.MethodPut).
		Path("/{stepId}/done").
		HandlerFunc(markInstructionStepAsDone)
	instructionStepRouter.
		Methods(http.MethodDelete).
		Path("/{stepId}/done").
		HandlerFunc(markInstructionStepAsTodo)

	apiRouter.
		Methods(http.MethodGet).
		Path("/healthz").
		HandlerFunc(getHealth)

	apiRouter.Use(authentication(), contentTypeJson)
}
