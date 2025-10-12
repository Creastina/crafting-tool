import { CloseEvent } from '../close-event.js';
import { get, put } from '../../lib/jinya-http.js';

class ProjectUpdatedEvent extends Event {
  constructor(project) {
    super('project-updated', {
      bubbles: true,
      cancelable: false,
      composed: true,
    });
    this.project = project;
  }
}

Alpine.data('updateProject', (project) => ({
  hasError: false,
  name: project.name,
  note: project.note,
  inventoryItems: project.inventoryItems,
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
      this.inventoryItems.push(item);
      this.foundItems = this.foundItems.filter((item) => item.id !== id);
      this.inventorySearch = '';
    }
  },
  deselectItem(id) {
    this.inventoryItems = this.inventoryItems.filter((item) => item.id !== id);
  },
}));

class UpdateProjectElement extends HTMLElement {
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
      <div class="creastina-dialog__container" id="dialog" x-data='updateProject(${this.getAttribute('project')})'>
        <div class="creastina-dialog">
          <header class="creastina-dialog__header">
            <h1 class="creastina-dialog__title">${this.getAttribute('name')} bearbeiten</h1>
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
                <li>
                  <span x-text="item.name"></span>
                  <button class="creastina-button is--negative is--icon" type="button" @click="deselectItem(item.id)">
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
                  <span x-text="item.note"></span>
                </li>
              </template>
            </ul>
            <div class="creastina-dialog__buttons">
              <button id="close" type="button" class="creastina-button">Nevermind</button>
              <button id="save" type="submit" class="creastina-button is--primary">Projekt speichern</button>
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

      data.inventoryItems = data.inventoryItems.map((item) => item.id);

      try {
        await put(`/api/project/category/${categoryId}/project/${this.getAttribute('project-id')}`, Alpine.raw(data));
        this.dispatchEvent(new ProjectUpdatedEvent(data));
      } catch (e) {
        projectData.hasError = true;
      }
    });

    Alpine.initTree(this.root);
  }
}

if (!customElements.get('creastina-update-project')) {
  customElements.define('creastina-update-project', UpdateProjectElement);
}

export async function updateProject(categoryId, project) {
  return new Promise((resolve) => {
    const container = document.createElement('creastina-update-project');
    container.setAttribute('category-id', categoryId);
    container.setAttribute('project-id', project.id);
    container.setAttribute('project', JSON.stringify(Alpine.raw(project)));
    document.body.appendChild(container);

    container.addEventListener('project-updated', (e) => {
      resolve(e.project);
      container.remove();
    });
    container.addEventListener('close', () => {
      resolve(false);
      container.remove();
    });
  });
}
