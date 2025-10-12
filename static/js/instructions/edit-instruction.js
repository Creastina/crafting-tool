import { CloseEvent } from '../close-event.js';
import { get, put } from '../../lib/jinya-http.js';

class InstructionUpdatedEvent extends Event {
  constructor(instruction) {
    super('instruction-updated', {
      bubbles: true,
      cancelable: false,
      composed: true,
    });
    this.instruction = instruction;
  }
}

Alpine.data('updateInstruction', (instruction, steps) => ({
  hasError: false,
  name: instruction.name,
  note: instruction.note,
  steps: steps.map((step) => ({
    ...step,
    uuid: crypto.randomUUID(),
  })),
  get data() {
    return {
      name: this.name,
      note: this.note,
      steps: this.steps,
    };
  },
  addStep(event, uuid) {
    this.steps.push({
      description: event.currentTarget.value,
      uuid: uuid,
      id: -1,
      done: false,
    });
    event.currentTarget.value = '';
    this.$nextTick(() => {
      this.$el.parentElement.querySelector(`[id='${uuid}']`).focus();
    });
  },
}));

class UpdateInstructionElement extends HTMLElement {
  constructor() {
    super();

    this.root = this.attachShadow({ mode: 'closed' });
  }

  connectedCallback() {
    this.root.innerHTML = `
      <style>
        @import '/static/css/button.css';
        @import '/static/css/form.css';
        @import '/static/css/dialog.css';
        @import '/static/css/message.css';
        @import '/static/css/gluten.css';
      </style>
      <div class="creastina-dialog__container" id="dialog" x-data='updateInstruction(${this.getAttribute('instruction')}, ${this.getAttribute('steps')})'>
        <div class="creastina-dialog">
          <header class="creastina-dialog__header">
            <h1 class="creastina-dialog__title">Anleitung bearbeiten</h1>
          </header>
          <form class="creastina-dialog__content">
            <div class="creastina-message is--negative" :class="{ 'is--hidden': !hasError }">
              <p>Die Anleitung gibt es schon</p>
            </div>
            <div class="creastina-form">
              <label for="name" class="creastina-form__label">Name</label>
              <input id="name" x-model="name" type="text" class="creastina-input" required>
              <label for="note" class="creastina-form__label">Notiz</label>
              <input id="note" x-model="note" type="text" class="creastina-input">
            </div>
            <b>Schritte</b>
            <div class="creastina-form is--one-column">
              <template x-for="step in steps" :key="step.uuid">
                <input :id="step.uuid" x-model="step.description" type="text" class="creastina-input">
              </template>
              <input type="text" class="creastina-input" @input="addStep($event, crypto.randomUUID())">
            </div>
            <div class="creastina-dialog__buttons">
              <button id="close" type="button" class="creastina-button">Nevermind</button>
              <button id="save" type="submit" class="creastina-button is--primary">Anleitung speichern</button>
            </div>
          </form>
        </div>
      </div>
    `;

    this.root.getElementById('close').addEventListener('click', () => {
      this.dispatchEvent(new CloseEvent());
    });
    this.root.querySelector('form').addEventListener('submit', async (e) => {
      e.preventDefault();

      const instructionData = Alpine.$data(this.root.querySelector('#dialog'));
      const data = instructionData.data;
      const steps = data.steps.map((step) => ({
        id: step.id,
        done: step.done,
        description: step.description,
      }));

      try {
        await put(`/api/instruction/${this.getAttribute('instruction-id')}`, Alpine.raw(data));
        await put(`/api/instruction/${this.getAttribute('instruction-id')}/steps`, Alpine.raw(steps));
        this.dispatchEvent(new InstructionUpdatedEvent(data));
      } catch (e) {
        instructionData.hasError = true;
      }
    });

    Alpine.initTree(this.root);
  }
}

if (!customElements.get('creastina-update-instruction')) {
  customElements.define('creastina-update-instruction', UpdateInstructionElement);
}

export async function updateInstruction(instruction, instructionId) {
  return new Promise(async (resolve) => {
    const steps = await get(`/api/instruction/${instructionId}/step`);
    const container = document.createElement('creastina-update-instruction');
    container.setAttribute('instruction-id', instructionId);
    container.setAttribute('instruction', JSON.stringify(instruction));
    container.setAttribute('steps', JSON.stringify(steps));
    document.body.appendChild(container);

    container.addEventListener('instruction-updated', (e) => {
      resolve(e.instruction);
      container.remove();
    });
    container.addEventListener('close', () => {
      resolve(false);
      container.remove();
    });
  });
}
