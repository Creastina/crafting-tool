import { CloseEvent } from '../close-event.js';
import { post } from '../../lib/jinya-http.js';

class InstructionCreatedEvent extends Event {
  constructor(instruction) {
    super('instruction-created', {
      bubbles: true,
      cancelable: false,
      composed: true,
    });
    this.instruction = instruction;
  }
}

Alpine.data('createInstruction', () => ({
  hasError: false,
  name: '',
  note: '',
  steps: [],
  get data() {
    return {
      name: this.name,
      note: this.note,
      steps: this.steps,
    };
  },
  addStep(event, uuid) {
    this.steps.push({ step: event.currentTarget.value, uuid: uuid });
    event.currentTarget.value = '';
    this.$nextTick(() => {
      this.$el.parentElement.querySelector(`[id='${uuid}']`).focus();
    });
  },
}));

class CreateInstructionElement extends HTMLElement {
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
      <div class="creastina-dialog__container" id="dialog" x-data="createInstruction">
        <div class="creastina-dialog">
          <header class="creastina-dialog__header">
            <h1 class="creastina-dialog__title">Neue Anleitung</h1>
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
                <input :id="step.uuid" x-model="step.step" type="text" class="creastina-input">
              </template>
              <input type="text" class="creastina-input" @input="addStep($event, crypto.randomUUID())">
            </div>
            <div class="creastina-dialog__buttons">
              <button id="close" type="button" class="creastina-button">Nevermind</button>
              <button id="save" type="submit" class="creastina-button is--primary">Anleitung anlegen</button>
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
      data.steps = data.steps.filter((step) => step.step !== '').map((step) => step.step);

      try {
        await post(`/api/instruction`, Alpine.raw(data));
        this.dispatchEvent(new InstructionCreatedEvent(data));
      } catch (e) {
        instructionData.hasError = true;
      }
    });

    Alpine.initTree(this.root);
  }
}

if (!customElements.get('creastina-create-instruction')) {
  customElements.define('creastina-create-instruction', CreateInstructionElement);
}

export async function createInstruction() {
  return new Promise((resolve) => {
    const container = document.createElement('creastina-create-instruction');
    document.body.appendChild(container);

    container.addEventListener('instruction-created', (e) => {
      resolve(e.instruction);
      container.remove();
    });
    container.addEventListener('close', () => {
      resolve(false);
      container.remove();
    });
  });
}
