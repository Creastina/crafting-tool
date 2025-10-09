import { CloseEvent } from '../closeEvent.js';
import { post } from '../../lib/jinya-http.js';

class BoxCreatedEvent extends Event {
  constructor(name) {
    super('box-created', {
      bubbles: true,
      cancelable: false,
      composed: true,
    });
    this.name = name;
  }
}

class CreateBoxElement extends HTMLElement {
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
      <div class="creastina-dialog__container">
        <div class="creastina-dialog">
          <header class="creastina-dialog__header">
            <h1 class="creastina-dialog__title">Neue Box</h1>
          </header>
          <form class="creastina-dialog__content">
            <div id="message" class="creastina-message is--negative is--hidden">
              <p>Der Name ist schon vergeben, bitte nimm einen anderen.</p>
            </div>
            <div class="creastina-form">
              <label for="name" class="creastina-form__label">Name</label>
              <input id="name" type="text" class="creastina-input" required>
            </div>
            <div class="creastina-dialog__buttons">
              <button id="close" type="button" class="creastina-button">Nevermind</button>
              <button id="save" type="submit" class="creastina-button is--primary">Box erstellen</button>
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
      const name = this.root.getElementById('name').value;

      try {
        await post('/api/inventory/box', { name });
        this.dispatchEvent(new BoxCreatedEvent(name));
      } catch (e) {
        this.root.getElementById('message').classList.remove('is--hidden');
      }
    });

    Alpine.initTree(this.root);
  }
}

if (!customElements.get('creastina-create-box')) {
  customElements.define('creastina-create-box', CreateBoxElement);
}

export async function createBox() {
  return new Promise((resolve) => {
    const container = document.createElement('creastina-create-box');
    document.body.appendChild(container);

    container.addEventListener('box-created', (e) => {
      resolve(e.name);
      container.remove();
    });
    container.addEventListener('close', () => {
      resolve(false);
      container.remove();
    });
  });
}
