export class CloseEvent extends Event {
  constructor() {
    super('close', {
      bubbles: true,
      cancelable: false,
      composed: true,
    });
  }
}
