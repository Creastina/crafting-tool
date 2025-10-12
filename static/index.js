import { setup } from './lib/jinya-alpine-tools.js';

document.addEventListener('DOMContentLoaded', async () => {
  await setup({
    defaultPage: 'inventory',
    baseScriptPath: '/static/js',
    routerBasePath: '/',
    openIdClientId: window.creastinaConfig.openIdClientId,
    openIdUrl: window.creastinaConfig.openIdUrl,
    openIdCallbackUrl: window.creastinaConfig.openIdCallbackUrl,
  });
});
