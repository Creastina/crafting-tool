import { get, httpDelete, put } from '../lib/jinya-http.js';
import confirm from '../lib/ui/confirm.js';
import alert from '../lib/ui/alert.js';
import { createCategory } from './projects/create-category.js';
import { createProject } from './projects/create-project.js';
import { updateProject } from './projects/edit-project.js';

Alpine.data('projectsData', () => ({
  categories: [],
  projects: [],
  selectedCategory: null,
  filteredProjects: [],
  loading: true,
  archived: false,
  searchBox() {
    const lowerQuery = Alpine.store('search').query.toLowerCase();
    this.filteredProjects = this.projects.filter(
      (item) =>
        (item.name.toLowerCase().includes(lowerQuery) || item.note.toLowerCase().includes(lowerQuery)) &&
        item.isArchived === this.archived,
    );
  },
  async init() {
    await this.loadCategories();
    if (this.categories.length > 0) {
      await this.selectCategory(this.categories[0]);
    }
    Alpine.store('search').setSearch(this.searchBox.bind(this));
    this.loading = false;
  },
  async loadCategories() {
    this.categories = await get('/api/project/category');
  },
  async selectCategory(category) {
    this.selectedCategory = category;
    this.projects = await get(`/api/project/category/${category.id}/project`);
    this.filteredProjects = this.projects.filter((item) => item.isArchived === this.archived);
  },
  async createCategory() {
    const newCategoryName = await createCategory();
    if (newCategoryName) {
      await this.loadCategories();
      const newCategory = this.categories.find((cat) => cat.name === newCategoryName);
      await this.selectCategory(newCategory);
    }
  },
  async createProject() {
    await createProject(this.selectedCategory.id);
    await this.selectCategory(this.selectedCategory);
  },
  async updateProject(project) {
    await updateProject(this.selectedCategory.id, project);
    await this.selectCategory(this.selectedCategory);
  },
  async deleteProject(project) {
    if (
      await confirm({
        title: 'Projekt löschen',
        message: `Möchtest du das Projekt "${project.name}" löschen? Wenn du das Projekt löscht wird es aus der Liste entfernt und die verwendeten Dinge im Inventar sind wieder verfügbar. Wenn du das Projekt als abgeschlossen markieren willst, archiviere es lieber. Dann wird auch das Inventar reduziert.`,
        declineLabel: 'Nevermind',
        approveLabel: 'Projekt löschen',
        negative: true,
      })
    ) {
      try {
        await httpDelete(`/api/project/category/${this.selectedCategory.id}/project/${project.id}`);
        await this.selectCategory(this.selectedCategory);
      } catch (e) {
        alert({
          title: 'Fehler beim Löschen',
          message: `Beim Löschen des Projekt "${project.name}" ist ein Fehler aufgetreten.`,
          closeLabel: 'Verdammt',
          negative: true,
        });
      }
    }
  },
  async archiveProject(project) {
    if (
      await confirm({
        title: 'Projekt archivieren',
        message: `Möchtest du das Projekt "${project.name}" archivieren? Wenn du das Projekt archivierst wird es aus der Liste entfernt. Außerdem wird das Inventar um die verbrauchte Anzahl reduziert.`,
        declineLabel: 'Nevermind',
        approveLabel: 'Projekt archivieren',
      })
    ) {
      try {
        await put(`/api/project/category/${this.selectedCategory.id}/project/${project.id}/archive`);
        await this.selectCategory(this.selectedCategory);
      } catch (e) {
        alert({
          title: 'Fehler beim Archivieren',
          message: `Beim Archivieren des Projekt "${project.name}" ist ein Fehler aufgetreten.`,
          closeLabel: 'Verdammt',
          negative: true,
        });
      }
    }
  },
  async unarchiveProject(project) {
    await httpDelete(`/api/project/category/${this.selectedCategory.id}/project/${project.id}/archive`);
    await this.selectCategory(this.selectedCategory);
  },
  async toggleArchived() {
    this.archived = !this.archived;
    this.searchBox();
  },
  async archiveAll() {
    if (
      await confirm({
        title: 'Alle Projekte archivieren',
        message: `Möchtest du alle Projekte archivieren? Das Inventar um die verbrauchte Anzahl reduziert.`,
        declineLabel: 'Nevermind',
        approveLabel: 'Projekte archivieren',
      })
    ) {
      try {
        await put(`/api/project/category/${this.selectedCategory.id}/archive`);
        await this.selectCategory(this.selectedCategory);
      } catch (e) {
        alert({
          title: 'Fehler beim Archivieren',
          message: `Leider konnten die Projekte nicht archiviert werden.`,
          closeLabel: 'Verdammt',
          negative: true,
        });
      }
    }
  },
}));
