import { CloseEvent } from '../closeEvent.js';
import { post } from '../../lib/jinya-http.js';

class ItemCreatedEvent extends Event {
  constructor(item) {
    super('item-created', {
      bubbles: true,
      cancelable: false,
      composed: true,
    });
    this.item = item;
  }
}

Alpine.data('createInventoryItem', () => ({
  hasError: false,
  name: '',
  note: '',
  unit: '',
  count: 0,
  properties: [],
  errorMessage: 'Die Sache gibt es schon',
  get data() {
    return {
      name: this.name,
      note: this.note,
      unit: this.unit,
      count: this.count,
      properties: this.properties,
    };
  },
  addProperty(event, uuid) {
    this.properties.push({ key: event.currentTarget.value, value: '', uuid: uuid });
    event.currentTarget.value = '';
    this.$nextTick(() => {
      this.$el.parentElement.querySelector(`[id='${uuid}']`).focus();
    });
  },
}));

class CreateItemElement extends HTMLElement {
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
      <div class="creastina-dialog__container" id="dialog" x-data="createInventoryItem">
        <div class="creastina-dialog">
          <header class="creastina-dialog__header">
            <h1 class="creastina-dialog__title">Neues ins Inventar</h1>
          </header>
          <form class="creastina-dialog__content">
            <div id="message" class="creastina-message is--negative" :class="{ 'is--hidden': !hasError }">
              <p x-text="errorMessage"></p>
            </div>
            <div class="creastina-form">
              <label for="name" class="creastina-form__label">Sache</label>
              <input id="name" x-model="name" type="text" class="creastina-input" required>
              <label for="note" class="creastina-form__label">Notiz</label>
              <input id="note" x-model="note" type="text" class="creastina-input">
              <label for="count" class="creastina-form__label">Anzahl</label>
              <input id="count" x-model.number="count" type="number" class="creastina-input" required>
              <label for="unit" class="creastina-form__label">Einheit</label>
              <input id="unit" x-model="unit" type="text" class="creastina-input" required>
            </div>
            <b>Eigenschaften</b>
            <div class="creastina-form">
              <template x-for="property in properties" :key="property.uuid">
                <div style="display: contents">
                  <input :id="property.uuid" x-model="property.key" type="text" class="creastina-input">
                  <input x-model="property.value" type="text" class="creastina-input">
                </div>
              </template>
              <input type="text" class="creastina-input" @input="addProperty($event, crypto.randomUUID())">
              <input type="text" class="creastina-input" disabled>
            </div>
            <div class="creastina-dialog__buttons">
              <button id="close" type="button" class="creastina-button">Nevermind</button>
              <button id="save" type="submit" class="creastina-button is--primary">Ins Inventar legen</button>
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

      const boxId = this.getAttribute('box-id');
      const inventoryData = Alpine.$data(this.root.querySelector('#dialog'));
      const data = inventoryData.data;
      if (data.properties.some((property) => property.key !== '' && property.value === '')) {
        inventoryData.hasError = true;
        inventoryData.errorMessage = 'Du hast leere Eigenschaften';
        return;
      }
      const allPropertyNames = data.properties.map((property) => property.key);
      if (allPropertyNames.length !== new Set(allPropertyNames).size) {
        inventoryData.hasError = true;
        inventoryData.errorMessage = 'Du hast doppelte Eigenschaften';
        return;
      }

      data.properties = data.properties.filter((property) => property.key !== '');

      try {
        await post(`/api/inventory/box/${boxId}/item`, Alpine.raw(data));
        this.dispatchEvent(new ItemCreatedEvent(name));
      } catch (e) {
        inventoryData.hasError = true;
      }
    });

    Alpine.initTree(this.root);
  }
}

if (!customElements.get('creastina-create-inventory-item')) {
  customElements.define('creastina-create-inventory-item', CreateItemElement);
}

export async function createInventoryItem(boxId) {
  return new Promise((resolve) => {
    const container = document.createElement('creastina-create-inventory-item');
    container.setAttribute('box-id', boxId);
    document.body.appendChild(container);

    container.addEventListener('item-created', (e) => {
      resolve(e.item);
      container.remove();
    });
    container.addEventListener('close', () => {
      resolve(false);
      container.remove();
    });
  });
}
