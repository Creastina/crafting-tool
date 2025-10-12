import { get, httpDelete } from '../lib/jinya-http.js';
import { createInstruction } from './instructions/create-instruction.js';
import confirm from '../lib/ui/confirm.js';
import alert from '../lib/ui/alert.js';
import { updateInstruction } from './instructions/edit-instruction.js';
import { instructionSteps } from './instructions/steps.js';

Alpine.data('instructionsData', () => ({
  instructions: [],
  filteredInstructions: [],
  loading: true,
  async init() {
    await this.loadInstructions();
    Alpine.store('search').setSearch(this.searchBox.bind(this));
    this.loading = false;
  },
  async loadInstructions() {
    this.instructions = await get('/api/instruction');
    this.filteredInstructions = this.instructions;
  },
  searchBox() {
    const lowerQuery = Alpine.store('search').query.toLowerCase();
    this.filteredInstructions = this.instructions.filter(
      (item) => item.name.toLowerCase().includes(lowerQuery) || item.note.toLowerCase().includes(lowerQuery),
    );
  },
  async createInstruction() {
    await createInstruction();
    await this.loadInstructions();
  },
  async showSteps(instruction) {
    await instructionSteps(instruction.id, instruction.name);
    await this.loadInstructions();
  },
  async editInstruction(instruction) {
    await updateInstruction(instruction, instruction.id);
    await this.loadInstructions();
  },
  async deleteInstruction(instruction) {
    if (
      await confirm({
        title: 'Anleitung löschen',
        message: `Möchtest du die Anleitung "${instruction.name}" löschen?`,
        declineLabel: 'Nevermind',
        approveLabel: 'Anleitung löschen',
        negative: true,
      })
    ) {
      try {
        await httpDelete(`/api/instruction/${instruction.id}`);
        await this.loadInstructions();
      } catch (e) {
        alert({
          title: 'Fehler beim Löschen',
          message: `Beim Löschen der Anleitung "${instruction.name}" ist ein Fehler aufgetreten.`,
          closeLabel: 'Verdammt',
          negative: true,
        });
      }
    }
  },
}));
