/**
 * Displays a confirm modal dialog
 * @param title {string}
 * @param message {string}
 * @param closeLabel {string}
 * @param negative {boolean}
 * @param positive {boolean}
 * @return {Promise<void>}
 */
export default async function alert({ title = window.location.href, message, closeLabel, negative = false }) {
  return new Promise((resolve) => {
    const container = document.createElement('div');
    document.body.appendChild(container);

    container.innerHTML = `
      <div class="creastina-dialog__container">
        <div class="creastina-dialog" style="max-width: 40rem">
          <header class="creastina-dialog__header ${negative ? 'is--negative' : ''}">
            <h1 class="creastina-dialog__title">${title}</h1>
          </header>
          <div class="creastina-dialog__content">
            <p class="creastina-dialog__message">${message}</p>
            <div class="creastina-dialog__buttons">
              <button data-action="close" class="creastina-button ${negative ? 'is--negative' : ''}">${closeLabel}</button>
            </div>
          </div>
        </div>
      </div>`;

    document.body.appendChild(container);

    container.querySelector('[data-action=close]').addEventListener('click', (e) => {
      e.preventDefault();
      container.remove();
      resolve();
    });
  });
}
