import { get, httpDelete, put } from '../lib/jinya-http.js';
import { createBox } from './inventory/create-box.js';
import { createInventoryItem } from './inventory/create-inventory.js';
import confirm from '../lib/ui/confirm.js';
import alert from '../lib/ui/alert.js';
import { updateInventoryItem } from './inventory/edit-inventory.js';

Alpine.data('inventoryData', () => ({
  boxes: [],
  boxContent: [],
  selectedBox: null,
  filteredBoxContent: [],
  loading: true,
  getProjectCount(count) {
    return `In ${count} ${count === 1 ? 'Projekt' : 'Projekten'} benötigt`;
  },
  getStockText(item) {
    return `${item.count} ${item.unit}`;
  },
  searchBox() {
    const lowerQuery = Alpine.store('search').query.toLowerCase();
    this.filteredBoxContent = this.boxContent.filter(
      (item) =>
        item.name.toLowerCase().includes(lowerQuery) ||
        item.note.toLowerCase().includes(lowerQuery) ||
        Object.values(item.properties).some((value) => value.toLowerCase().includes(lowerQuery)),
    );
  },
  async init() {
    await this.loadCategories();
    if (this.boxes.length > 0) {
      await this.selectBox(this.boxes[0]);
    }
    Alpine.store('search').setSearch(this.searchBox.bind(this));
    this.loading = false;
  },
  async loadCategories() {
    this.boxes = await get('/api/inventory/box');
  },
  async selectBox(box) {
    this.selectedBox = box;
    this.boxContent = await get(`/api/inventory/box/${box.id}/item`);
    this.filteredBoxContent = this.boxContent;
  },
  async createProject() {
    const newBoxName = await createBox();
    if (newBoxName) {
      await this.loadCategories();
      const newBox = this.boxes.find((box) => box.name === newBoxName);
      await this.selectBox(newBox);
    }
  },
  async createItem() {
    await createInventoryItem(this.selectedBox.id);
    await this.selectBox(this.selectedBox);
  },
  async updateItem(item) {
    await updateInventoryItem(this.selectedBox.id, item);
    await this.selectBox(this.selectedBox);
  },
  async deleteItem(item) {
    if (
      await confirm({
        title: 'Aus dem Inventar entfernen',
        message: `Möchtest du "${item.name}" aus dem Inventar entfernen? Wenn du "${item.name}" entfernst, kannst es nicht mehr in deinen Projekten verwenden.`,
        declineLabel: 'Nevermind',
        approveLabel: 'Rausnehmen',
        negative: true,
      })
    ) {
      try {
        await httpDelete(`/api/inventory/box/${this.selectedBox.id}/item/${item.id}`);
        await this.selectBox(this.selectedBox);
      } catch (e) {
        alert({
          title: 'Fehler beim Löschen',
          message: `Beim Löschen der Sache "${item.name}" ist ein Fehler aufgetreten.`,
          closeLabel: 'Verdammt',
          negative: true,
        });
      }
    }
  },
  async increaseStock(item) {
    try {
      await put(`/api/inventory/box/${this.selectedBox.id}/item/${item.id}/stock`);
      await this.selectBox(this.selectedBox);
    } catch (e) {
      alert({
        title: 'Fehler beim Erhöhen',
        message: `Beim Erhöhen des Inventars von "${item.name}"`,
        closeLabel: 'Verdammt',
        negative: true,
      });
    }
  },
  async decreaseStock(item) {
    try {
      await httpDelete(`/api/inventory/box/${this.selectedBox.id}/item/${item.id}/stock`);
      await this.selectBox(this.selectedBox);
    } catch (e) {
      alert({
        title: 'Fehler beim Verringern',
        message: `Beim Verringern des Inventars von "${item.name}"`,
        closeLabel: 'Verdammt',
        negative: true,
      });
    }
  },
}));
