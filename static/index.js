import { setup } from './lib/jinya-alpine-tools.js';

document.addEventListener('DOMContentLoaded', async () => {
  await setup({
    defaultPage: 'inventory',
    baseScriptPath: '/static/js/',
    routerBasePath: '/',
    openIdClientId: window.jewelsConfig.openIdClientId,
    openIdUrl: window.jewelsConfig.openIdUrl,
    openIdCallbackUrl: window.jewelsConfig.openIdCallbackUrl,
  });
});
