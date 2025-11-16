import { CloseEvent } from '../close-event.js';
import { get, post } from '../../lib/jinya-http.js';
import { setItemCount } from './count-setter.js';

class ProjectCreatedEvent extends Event {
  constructor(project) {
    super('project-created', {
      bubbles: true,
      cancelable: false,
      composed: true,
    });
    this.project = project;
  }
}

Alpine.data('createProject', () => ({
  hasError: false,
  name: '',
  note: '',
  inventoryItems: [],
  inventorySearch: '',
  foundItems: [],
  errorMessage: 'Das Projekt gibt es schon',
  get data() {
    return {
      name: this.name,
      note: this.note,
      inventoryItems: this.inventoryItems,
    };
  },
  init() {
    this.$watch('inventorySearch', async (value) => {
      if (this.inventorySearch.length < 3) {
        this.foundItems = [];
        return;
      }

      try {
        const data = await get(`/api/inventory/search?query=${value}`);
        this.foundItems = data.filter(
          (item) => !this.inventoryItems.some((inventoryItem) => inventoryItem.id === item.id),
        );
      } catch (e) {
        console.error(e);
      }
    });
  },
  selectItem(id) {
    const item = this.foundItems.find((item) => item.id === id);
    if (item) {
      item.count = 1;
      this.inventoryItems.push(item);
      this.foundItems = this.foundItems.filter((item) => item.id !== id);
      this.inventorySearch = '';
    }
  },
  deselectItem(id) {
    this.inventoryItems = this.inventoryItems.filter((item) => item.id !== id);
  },
  async changeCount(item) {
    const count = await setItemCount(item, item.count);
    const itemIdx = this.inventoryItems.indexOf(item);
    this.inventoryItems[itemIdx].count = count;
  },
}));

class CreateProjectElement extends HTMLElement {
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
      <div class="creastina-dialog__container" id="dialog" x-data="createProject">
        <div class="creastina-dialog">
          <header class="creastina-dialog__header">
            <h1 class="creastina-dialog__title">Neues Projekt</h1>
          </header>
          <form class="creastina-dialog__content">
            <div id="message" class="creastina-message is--negative" :class="{ 'is--hidden': !hasError }">
              <p x-text="errorMessage"></p>
            </div>
            <div class="creastina-form">
              <label for="name" class="creastina-form__label">Projekt</label>
              <input id="name" x-model="name" type="text" class="creastina-input" required>
              <label for="note" class="creastina-form__label">Notiz</label>
              <input id="note" x-model="note" type="text" class="creastina-input">
            </div>
            <b>Benötigtes Inventar</b>
            <div class="creastina-form">
              <label for="search" class="creastina-form__label">Suche</label>
              <input type="search" class="creastina-input" x-model="inventorySearch" x-ref="search">
            </div>
            <div class="creastina-search__list" x-anchor="$refs.search" x-collapse x-show="inventorySearch.length >= 3">
              <template x-for="item in foundItems" :key="item.id">
                <div class="creastina-search__item">
                  <button class="creastina-search__item is--button" type="button" x-text="item.name" @click="selectItem(item.id)"></button>
                </div>
              </template>
            </div>
            <ul class="creastina-project__inventory">
              <template x-for="item in inventoryItems" :key="item.id">
                <li class="creastina-project__inventory-item">
                  <span class="creastina-project__inventory-name" x-text="item.name"></span>
                  <span class="creastina-project__inventory-note" x-text="item.note"></span>
                  <button class="creastina-button is--accent-1 is--icon is--item-count" type="button" @click="changeCount(item)"">
                    <span x-text="item.count"></span>
                    <svg viewBox="0 0 24 24">
                      <path
                        fill="none"
                        stroke="currentColor"
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        stroke-width="2"
                        d="M3.5 18.985V20.5h1.514c1.227 0 1.84 0 2.391-.228c.551-.229.985-.662 1.852-1.53l9.864-9.863c.883-.883 1.324-1.324 1.373-1.866q.012-.135 0-.269c-.05-.541-.49-.983-1.373-1.865c-.883-.883-1.324-1.324-1.865-1.373a1.5 1.5 0 0 0-.27 0c-.541.049-.982.49-1.865 1.373l-9.864 9.864c-.867.867-1.3 1.3-1.529 1.852c-.228.55-.228 1.164-.228 2.39M13.5 6.5l4 4"
                      />
                    </svg>
                  </button>
                  <button class="creastina-button is--negative is--icon is--delete-item" type="button" @click="deselectItem(item.id)">
                    <svg viewBox="0 0 24 24">
                      <path
                        fill="none"
                        stroke="currentColor"
                        stroke-linecap="round"
                        stroke-width="2"
                        d="m19.5 5.5l-.62 10.025c-.158 2.561-.237 3.842-.88 4.763a4 4 0 0 1-1.2 1.128c-.957.584-2.24.584-4.806.584c-2.57 0-3.855 0-4.814-.585a4 4 0 0 1-1.2-1.13c-.642-.922-.72-2.205-.874-4.77L4.5 5.5M9 11.735h6m-4.5 3.919h3M3 5.5h18m-4.945 0l-.682-1.408c-.454-.936-.68-1.403-1.071-1.695a2 2 0 0 0-.275-.172C13.594 2 13.074 2 12.034 2c-1.065 0-1.598 0-2.039.234a2 2 0 0 0-.278.18c-.396.303-.617.788-1.059 1.757L8.053 5.5"
                      />
                    </svg>
                  </button>
                </li>
              </template>
            </ul>
            <div class="creastina-dialog__buttons">
              <button id="close" type="button" class="creastina-button">Nevermind</button>
              <button id="save" type="submit" class="creastina-button is--primary">Projekt erstellen</button>
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

      const categoryId = this.getAttribute('category-id');
      const projectData = Alpine.$data(this.root.querySelector('#dialog'));
      const data = projectData.data;

      data.inventoryItems = data.inventoryItems.reduce((acc, item) => {
        acc[item.id] = item.count;
        return acc;
      }, {});

      try {
        await post(`/api/project/category/${categoryId}/project`, Alpine.raw(data));
        this.dispatchEvent(new ProjectCreatedEvent(data));
      } catch (e) {
        projectData.hasError = true;
      }
    });

    Alpine.initTree(this.root);
  }
}

if (!customElements.get('creastina-create-project')) {
  customElements.define('creastina-create-project', CreateProjectElement);
}

export async function createProject(categoryId) {
  return new Promise((resolve) => {
    const container = document.createElement('creastina-create-project');
    container.setAttribute('category-id', categoryId);
    document.body.appendChild(container);

    container.addEventListener('project-created', (e) => {
      resolve(e.project);
      container.remove();
    });
    container.addEventListener('close', () => {
      resolve(false);
      container.remove();
    });
  });
}
