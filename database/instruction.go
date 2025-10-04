package database

func GetInstructions() ([]Instruction, error) {
	return Select[Instruction]("select * from instruction")
}

func GetInstruction(id int) (*Instruction, error) {
	return Get[Instruction](id)
}

func CreateInstruction(instruction Instruction, steps []string) (*Instruction, error) {
	tx, err := dbMap.Begin()
	if err != nil {
		return nil, err
	}

	err = tx.Insert(&instruction)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	for _, step := range steps {
		instructionStep := InstructionStep{
			InstructionId: instruction.Id,
			Description:   step,
			Done:          false,
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

	for _, step := range steps {
		_, err = tx.Exec(`
insert into instruction_step (instruction_id, description, done) values ($1, $2, $3)`, id, step.Description, step.Done)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func GetInstructionSteps(instructionId int) ([]InstructionStep, error) {
	return Select[InstructionStep]("select * from instruction_step where instruction_id = $1", instructionId)
}

func GetInstructionStep(stepId, instructionId int) (*InstructionStep, error) {
	step, err := SelectOne[InstructionStep]("select * from instruction_step where id = $1 instruction_id = $2", stepId, instructionId)
	if err != nil {
		return nil, err
	}

	return &step, nil
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
