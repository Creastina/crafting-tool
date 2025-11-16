import { CloseEvent } from '../close-event.js';

class CountSetEvent extends Event {
  constructor(count) {
    super('count-set', {
      bubbles: true,
      cancelable: false,
      composed: true,
    });
    this.count = count;
  }
}

Alpine.data('setItemCount', (itemName, count) => ({
  hasError: false,
  itemName,
  count,
  get data() {
    return {
      count: this.count,
    };
  },
}));

class SetItemCountElement extends HTMLElement {
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
        @import '/static/css/projects.css';
      </style>
      <div class="creastina-dialog__container" id="dialog" x-data="setItemCount('${this.getAttribute('item-name')}', ${this.getAttribute('count')})">
        <div class="creastina-dialog">
          <header class="creastina-dialog__header">
            <h1 class="creastina-dialog__title" x-text="'Benötigte Anzahl für ' + itemName"></h1>
          </header>
          <form class="creastina-dialog__content">
            <div class="creastina-form">
              <label for="count" class="creastina-form__label">Anzahl</label>
              <input id="count" x-model.number="count" type="text" class="creastina-input" required>
            </div>
            <div class="creastina-dialog__buttons">
              <button id="close" type="button" class="creastina-button">Nevermind</button>
              <button id="save" type="submit" class="creastina-button is--primary">Anzahl festlegen</button>
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

      const countData = Alpine.$data(this.root.querySelector('#dialog'));
      const data = countData.data;

      this.dispatchEvent(new CountSetEvent(data.count));
    });

    Alpine.initTree(this.root);
  }
}

if (!customElements.get('creastina-set-item-count')) {
  customElements.define('creastina-set-item-count', SetItemCountElement);
}

export async function setItemCount(item, count) {
  return new Promise((resolve) => {
    const container = document.createElement('creastina-set-item-count');
    container.setAttribute('item-name', item.name);
    container.setAttribute('count', count);
    document.body.appendChild(container);

    container.addEventListener('count-set', (e) => {
      resolve(e.count);
      container.remove();
    });
    container.addEventListener('close', () => {
      resolve(false);
      container.remove();
    });
  });
}
