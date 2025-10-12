import { CloseEvent } from '../close-event.js';
import { get, httpDelete, put } from '../../lib/jinya-http.js';

Alpine.data('instructionSteps', (name, id, steps) => ({
  hasError: false,
  steps: steps,
  async markDone(step) {
    await put(`/api/instruction/${id}/step/${step.id}/done`);
    step.done = true;
  },
  async markTodo(step) {
    await httpDelete(`/api/instruction/${id}/step/${step.id}/done`);
    step.done = false;
  },
}));

class InstructionStepsElement extends HTMLElement {
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
        @import '/static/css/steps.css';
        @import '/static/css/gluten.css';
      </style>
      <div class="creastina-dialog__container" id="dialog" x-data='instructionSteps("${this.getAttribute('instruction-name')}", ${this.getAttribute('instruction-id')}, ${this.getAttribute('steps')})'>
        <div class="creastina-dialog">
          <header class="creastina-dialog__header">
            <h1 class="creastina-dialog__title">Schritte für "${this.getAttribute('instruction-name')}"</h1>
          </header>
          <div class="creastina-dialog__content">
            <div class="creastina-steps">
              <template x-for="step in steps">
                <div style="display: contents">
                  <template x-if="step.done">
                    <button class="creastina-button is--step is--done" @click="markTodo(step)">
                      <svg viewBox="0 0 24 24">
                        <g fill="none" stroke="currentColor" stroke-linecap="round" stroke-width="2">
                          <path d="M21.448 8.2c.052 1.05.052 2.3.052 3.8c0 4.478 0 6.718-1.391 8.109S16.479 21.5 12 21.5c-4.478 0-6.718 0-8.109-1.391S2.5 16.479 2.5 12c0-4.478 0-6.718 1.391-8.109S7.521 2.5 12 2.5c1.072 0 2.016 0 2.85.02"/>
                          <path class="creastina-steps__checkmark" stroke-linejoin="round" d="M8 11.5s1.5 0 3.5 3.5c0 0 5.059-9.167 10-11"/>
                        </g>
                      </svg>
                      <span class="creastina-steps__step-text" x-text="step.description"></span>
                    </button>
                  </template>
                  <template x-if="!step.done">
                    <button class="creastina-button is--step is--todo" @click="markDone(step)">
                      <svg viewBox="0 0 24 24">
                        <g fill="none" stroke="currentColor" stroke-linecap="round" stroke-width="2">
                          <path d="M21.448 8.2c.052 1.05.052 2.3.052 3.8c0 4.478 0 6.718-1.391 8.109S16.479 21.5 12 21.5c-4.478 0-6.718 0-8.109-1.391S2.5 16.479 2.5 12c0-4.478 0-6.718 1.391-8.109S7.521 2.5 12 2.5c1.072 0 2.016 0 2.85.02"/>
                          <path class="creastina-steps__checkmark" stroke-linejoin="round" d="M8 11.5s1.5 0 3.5 3.5c0 0 5.059-9.167 10-11"/>
                        </g>
                      </svg>
                      <span class="creastina-steps__step-text" x-text="step.description"></span>
                    </button>
                    <span class="creastina-steps__step-text" x-text="step.description"></span>
                  </template>
                </div>
              </template>
            </div>
            <div class="creastina-dialog__buttons">
              <button id="done" class="creastina-button is--primary">Fertig</button>
            </div>
          </d>
        </div>
      </div>
    `;

    this.root.getElementById('done').addEventListener('click', () => {
      this.dispatchEvent(new CloseEvent());
    });

    Alpine.initTree(this.root);
  }
}

if (!customElements.get('creastina-instruction-steps')) {
  customElements.define('creastina-instruction-steps', InstructionStepsElement);
}

export async function instructionSteps(instructionId, instructionName) {
  return new Promise(async (resolve) => {
    const steps = await get(`/api/instruction/${instructionId}/step`);
    const container = document.createElement('creastina-instruction-steps');
    container.setAttribute('instruction-id', instructionId);
    container.setAttribute('instruction-name', instructionName);
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
