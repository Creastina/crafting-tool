package api

import (
	"crafting/database"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func getInventoryBoxes(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)

	boxes, err := database.Select[database.InventoryBox]("select * from inventory_box order by name")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = encoder.Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	_ = encoder.Encode(boxes)
}

func getInventoryBox(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)

	vars := mux.Vars(r)

	boxId, err := strconv.Atoi(vars["boxId"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = encoder.Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	box, err := database.Get[database.InventoryBox](boxId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = encoder.Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	_ = encoder.Encode(box)
}

func createInventoryBox(w http.ResponseWriter, r *http.Request) {
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

	box := database.InventoryBox{
		Name: body.Name,
	}

	err := database.GetDbMap().Insert(&box)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = encoder.Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusCreated)
	_ = encoder.Encode(box)
}

func getInventoryItems(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)

	vars := mux.Vars(r)

	boxId, err := strconv.Atoi(vars["boxId"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = encoder.Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	items, err := database.GetInventoryItems(boxId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = encoder.Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	_ = encoder.Encode(items)
}

func getInventoryItem(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)

	vars := mux.Vars(r)

	itemId, err := strconv.Atoi(vars["itemId"])
	boxId, err := strconv.Atoi(vars["boxId"])

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = encoder.Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	item, err := database.GetInventoryItem(itemId, boxId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = encoder.Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	_ = encoder.Encode(item)
}

func createInventoryItem(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	decoder := json.NewDecoder(r.Body)

	vars := mux.Vars(r)

	boxId, err := strconv.Atoi(vars["boxId"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = encoder.Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	body := struct {
		Name       string            `json:"name"`
		Note       string            `json:"note"`
		Count      int               `json:"count"`
		Unit       string            `json:"unit"`
		Properties map[string]string `json:"properties"`
	}{}
	if err := decoder.Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = encoder.Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	item := database.InventoryItem{
		Name:       body.Name,
		Note:       body.Note,
		Count:      body.Count,
		Unit:       body.Unit,
		BoxId:      boxId,
		Properties: body.Properties,
	}

	result, err := database.CreateInventoryItem(item)
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

func updateInventoryItem(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	decoder := json.NewDecoder(r.Body)

	vars := mux.Vars(r)

	boxId, err := strconv.Atoi(vars["boxId"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = encoder.Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	itemId, err := strconv.Atoi(vars["itemId"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = encoder.Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	body := struct {
		Name       string            `json:"name"`
		Note       string            `json:"note"`
		Count      int               `json:"count"`
		Unit       string            `json:"unit"`
		Properties map[string]string `json:"properties"`
	}{}
	if err := decoder.Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = encoder.Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	item := database.InventoryItem{
		Id:         itemId,
		Name:       body.Name,
		Note:       body.Note,
		Count:      body.Count,
		Unit:       body.Unit,
		BoxId:      boxId,
		Properties: body.Properties,
	}

	err = database.UpdateInventoryItem(item)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = encoder.Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func deleteInventoryItem(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)

	vars := mux.Vars(r)

	boxId, err := strconv.Atoi(vars["boxId"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = encoder.Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	itemId, err := strconv.Atoi(vars["itemId"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = encoder.Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	err = database.DeleteInventoryItem(itemId, boxId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = encoder.Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
