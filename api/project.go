package api

import (
	"crafting/database"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func getProjectCategories(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)

	projectCategories, err := database.Select[database.ProjectCategory]("select * from project_category order by name")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = encoder.Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	_ = encoder.Encode(projectCategories)
}

func getProjectCategory(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)

	vars := mux.Vars(r)
	categoryId, err := strconv.Atoi(vars["categoryId"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = encoder.Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	projectCategory, err := database.Get[database.ProjectCategory](categoryId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = encoder.Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	_ = encoder.Encode(projectCategory)
}

func createProjectCategory(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	decoder := json.NewDecoder(r.Body)

	body := struct {
		Name string `json:"name"`
	}{}
	if err := decoder.Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = encoder.Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	projectCategory := database.ProjectCategory{
		Name: body.Name,
	}

	err := database.GetDbMap().Insert(&projectCategory)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = encoder.Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	_ = encoder.Encode(projectCategory)
}

func getProjects(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)

	isArchived := r.URL.Query().Get("archived") == "true"

	vars := mux.Vars(r)
	categoryId, err := strconv.Atoi(vars["categoryId"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = encoder.Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	projects, err := database.GetProjects(categoryId, isArchived)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = encoder.Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	_ = encoder.Encode(projects)
}

func getProject(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)

	vars := mux.Vars(r)

	projectId, err := strconv.Atoi(vars["projectId"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = encoder.Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	categoryId, err := strconv.Atoi(vars["categoryId"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = encoder.Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	project, err := database.GetProject(projectId, categoryId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = encoder.Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	_ = encoder.Encode(project)
}

func createProject(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	decoder := json.NewDecoder(r.Body)

	vars := mux.Vars(r)
	categoryId, err := strconv.Atoi(vars["categoryId"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = encoder.Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	body := struct {
		Name           string `json:"name"`
		Note           string `json:"note"`
		InventoryItems []int  `json:"inventoryItems"`
	}{}
	if err := decoder.Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = encoder.Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	project := database.Project{
		Name:       body.Name,
		Note:       body.Note,
		CategoryId: categoryId,
	}

	result, err := database.CreateProject(project, body.InventoryItems)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = encoder.Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusCreated)
	_ = encoder.Encode(result)
}

func updateProject(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	decoder := json.NewDecoder(r.Body)
	vars := mux.Vars(r)

	categoryId, err := strconv.Atoi(vars["categoryId"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = encoder.Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	projectId, err := strconv.Atoi(vars["projectId"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = encoder.Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	body := struct {
		Name           string `json:"name"`
		Note           string `json:"note"`
		InventoryItems []int  `json:"inventoryItems"`
		IsArchived     bool   `json:"isArchived"`
	}{}

	if err := decoder.Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = encoder.Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	project := database.Project{
		Id:         projectId,
		CategoryId: categoryId,
		Name:       body.Name,
		Note:       body.Note,
		IsArchived: body.IsArchived,
	}

	err = database.UpdateProject(project, body.InventoryItems)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = encoder.Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func deleteProject(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)

	vars := mux.Vars(r)

	categoryId, err := strconv.Atoi(vars["categoryId"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = encoder.Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	projectId, err := strconv.Atoi(vars["projectId"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = encoder.Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	err = database.DeleteProject(projectId, categoryId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = encoder.Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func archiveProject(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)

	vars := mux.Vars(r)

	categoryId, err := strconv.Atoi(vars["categoryId"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = encoder.Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	projectId, err := strconv.Atoi(vars["projectId"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = encoder.Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	_, err = database.GetDbMap().Exec("update project set is_archived = true where id = $1 and category_id = $2", projectId, categoryId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = encoder.Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func unarchiveProject(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)

	vars := mux.Vars(r)

	categoryId, err := strconv.Atoi(vars["categoryId"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = encoder.Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	projectId, err := strconv.Atoi(vars["projectId"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = encoder.Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	_, err = database.GetDbMap().Exec("update project set is_archived = false where id = $1 and category_id = $2", projectId, categoryId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = encoder.Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
