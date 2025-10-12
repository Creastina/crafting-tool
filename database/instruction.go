package database

func GetInstructions() ([]InstructionWithStepCount, error) {
	return Select[InstructionWithStepCount]("select * from instruction_with_step_count")
}

func GetInstruction(id int) (*InstructionWithStepCount, error) {
	return SelectOne[InstructionWithStepCount]("select * from instruction_with_step_count where id = $1", id)
}

func CreateInstruction(instruction Instruction, steps []string) (*InstructionWithStepCount, error) {
	tx, err := dbMap.Begin()
	if err != nil {
		return nil, err
	}

	err = tx.Insert(&instruction)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	for i, step := range steps {
		instructionStep := InstructionStep{
			InstructionId: instruction.Id,
			Description:   step,
			Done:          false,
			Position:      i,
		}

		err = tx.Insert(&instructionStep)
		if err != nil {
			_ = tx.Rollback()
			return nil, err
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return GetInstruction(instruction.Id)
}

func UpdateInstruction(instruction Instruction) error {
	_, err := dbMap.Update(&instruction)

	return err
}

func DeleteInstruction(id int) error {
	_, err := dbMap.Exec("delete from instruction where id = $1", id)

	return err
}

func ReplaceInstructionSteps(id int, steps []InstructionStep) error {
	tx, err := dbMap.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec("delete from instruction_step where instruction_id = $1", id)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	for i, step := range steps {
		_, err = tx.Exec(`
insert into instruction_step (instruction_id, description, done, position) values ($1, $2, $3, $4)`, id, step.Description, step.Done, i)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func GetInstructionSteps(instructionId int) ([]InstructionStep, error) {
	return Select[InstructionStep]("select * from instruction_step where instruction_id = $1 order by id", instructionId)
}

func GetInstructionStep(stepId, instructionId int) (*InstructionStep, error) {
	return SelectOne[InstructionStep]("select * from instruction_step where id = $1 instruction_id = $2", stepId, instructionId)
}

func CreateInstructionStep(instructionStep InstructionStep) (*InstructionStep, error) {
	err := dbMap.Insert(&instructionStep)
	if err != nil {
		return nil, err
	}

	return &instructionStep, nil
}

func UpdateInstructionStep(instructionStep InstructionStep) error {
	_, err := dbMap.Update(&instructionStep)
	return err
}

func DeleteInstructionStep(stepId, instructionId int) error {
	_, err := dbMap.Exec("delete from instruction_step where id = $1 and instruction_id = $2", stepId, instructionId)
	return err
}

func MarkInstructionStepAsDone(stepId, instructionId int) error {
	_, err := dbMap.Exec("update instruction_step set done = true where id = $1 and instruction_id = $2", stepId, instructionId)
	return err
}

func MarkInstructionStepAsTodo(stepId, instructionId int) error {
	_, err := dbMap.Exec("update instruction_step set done = false where id = $1 and instruction_id = $2", stepId, instructionId)
	return err
}
