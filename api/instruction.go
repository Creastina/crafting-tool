package api

import (
	"crafting/database"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func getInstructions(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)

	instructions, err := database.GetInstructions()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = encoder.Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	_ = encoder.Encode(instructions)
}

func getInstruction(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)

	vars := mux.Vars(r)

	instructionId, err := strconv.Atoi(vars["instructionId"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = encoder.Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	instruction, err := database.GetInstruction(instructionId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = encoder.Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	_ = encoder.Encode(instruction)
}

func createInstruction(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	decoder := json.NewDecoder(r.Body)

	body := struct {
		Name  string   `json:"name"`
		Note  string   `json:"note"`
		Steps []string `json:"steps"`
	}{}

	err := decoder.Decode(&body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = encoder.Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	instruction := database.Instruction{
		Name: body.Name,
		Note: body.Note,
	}

	result, err := database.CreateInstruction(instruction, body.Steps)
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

func updateInstruction(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	decoder := json.NewDecoder(r.Body)

	vars := mux.Vars(r)
	instructionId, err := strconv.Atoi(vars["instructionId"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = encoder.Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	body := struct {
		Name string `json:"name"`
		Note string `json:"note"`
	}{}

	err = decoder.Decode(&body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = encoder.Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	instruction := database.Instruction{
		Id:   instructionId,
		Name: body.Name,
		Note: body.Note,
	}

	err = database.UpdateInstruction(instruction)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = encoder.Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func deleteInstruction(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)

	vars := mux.Vars(r)
	instructionId, err := strconv.Atoi(vars["instructionId"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = encoder.Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	err = database.DeleteInstruction(instructionId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = encoder.Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func replaceSteps(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	decoder := json.NewDecoder(r.Body)

	vars := mux.Vars(r)
	instructionId, err := strconv.Atoi(vars["instructionId"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = encoder.Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	var body []struct {
		Description string `json:"description"`
		Done        bool   `json:"done"`
	}

	err = decoder.Decode(&body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = encoder.Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	instructionSteps := make([]database.InstructionStep, len(body))
	for i, step := range body {
		instructionSteps[i] = database.InstructionStep{
			InstructionId: instructionId,
			Description:   step.Description,
			Done:          step.Done,
		}
	}

	err = database.ReplaceInstructionSteps(instructionId, instructionSteps)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = encoder.Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func getInstructionSteps(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)

	vars := mux.Vars(r)
	instructionId, err := strconv.Atoi(vars["instructionId"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = encoder.Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	steps, err := database.GetInstructionSteps(instructionId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = encoder.Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	_ = encoder.Encode(steps)
}

func getInstructionStep(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)

	vars := mux.Vars(r)
	instructionId, err := strconv.Atoi(vars["instructionId"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = encoder.Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	stepId, err := strconv.Atoi(vars["stepId"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = encoder.Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	steps, err := database.GetInstructionStep(stepId, instructionId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = encoder.Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	_ = encoder.Encode(steps)
}

func createInstructionStep(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	decoder := json.NewDecoder(r.Body)

	vars := mux.Vars(r)
	instructionId, err := strconv.Atoi(vars["instructionId"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = encoder.Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	body := struct {
		Description string `json:"description"`
	}{}

	err = decoder.Decode(&body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = encoder.Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	step := database.InstructionStep{
		InstructionId: instructionId,
		Description:   body.Description,
	}

	result, err := database.CreateInstructionStep(step)
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

func updateInstructionStep(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	decoder := json.NewDecoder(r.Body)

	vars := mux.Vars(r)
	instructionId, err := strconv.Atoi(vars["instructionId"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = encoder.Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	stepId, err := strconv.Atoi(vars["stepId"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = encoder.Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	body := struct {
		Description string `json:"description"`
		Done        bool   `json:"done"`
	}{}

	err = decoder.Decode(&body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = encoder.Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	step := database.InstructionStep{
		Id:            stepId,
		InstructionId: instructionId,
		Description:   body.Description,
		Done:          body.Done,
	}

	err = database.UpdateInstructionStep(step)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = encoder.Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func deleteInstructionStep(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	vars := mux.Vars(r)
	instructionId, err := strconv.Atoi(vars["instructionId"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = encoder.Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	stepId, err := strconv.Atoi(vars["stepId"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = encoder.Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	err = database.DeleteInstructionStep(stepId, instructionId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = encoder.Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func markInstructionStepAsDone(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	vars := mux.Vars(r)
	instructionId, err := strconv.Atoi(vars["instructionId"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = encoder.Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	stepId, err := strconv.Atoi(vars["stepId"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = encoder.Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	err = database.MarkInstructionStepAsDone(stepId, instructionId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = encoder.Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func markInstructionStepAsTodo(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	vars := mux.Vars(r)
	instructionId, err := strconv.Atoi(vars["instructionId"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = encoder.Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	stepId, err := strconv.Atoi(vars["stepId"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = encoder.Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	err = database.MarkInstructionStepAsTodo(stepId, instructionId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = encoder.Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
