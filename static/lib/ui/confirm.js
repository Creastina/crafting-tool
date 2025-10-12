/**
 * Displays a confirm modal dialog
 * @param title {string}
 * @param message {string}
 * @param declineLabel {string}
 * @param approveLabel {string}
 * @param negative {boolean}
 * @param warning {boolean}
 * @return {Promise<boolean>}
 */
export default async function confirm({
  title = window.location.href,
  message,
  declineLabel,
  approveLabel,
  negative = false,
  warning = false,
}) {
  return new Promise((resolve) => {
    const container = document.createElement('div');
    document.body.appendChild(container);

    container.innerHTML = `
      <div class="creastina-dialog__container">
        <div class="creastina-dialog" style="max-width: 40rem">
          <header class="creastina-dialog__header ${negative ? 'is--negative' : ''} ${warning ? 'is--warning' : ''}">
            <h1 class="creastina-dialog__title">${title}</h1>
          </header>
          <div class="creastina-dialog__content">
            <p class="creastina-dialog__message">${message}</p>
            <div class="creastina-dialog__buttons">
                <button data-action="decline" class="creastina-button">${declineLabel}</button>
                <button data-action="approve" class="creastina-button ${!(negative && warning) ? 'is--primary' : ''} ${negative ? 'is--negative' : ''} ${warning ? 'is--warning' : ''}">${approveLabel}</button>
            </div>
          </div>
        </div>
      </div>`;

    document.body.appendChild(container);

    container.querySelector('[data-action=decline]').addEventListener('click', (e) => {
      e.preventDefault();
      container.remove();
      resolve(false);
    });
    container.querySelector('[data-action=approve]').addEventListener('click', (e) => {
      e.preventDefault();
      container.remove();
      resolve(true);
    });
  });
}
