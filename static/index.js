import { setup } from './lib/jinya-alpine-tools.js';

document.addEventListener('DOMContentLoaded', async () => {
  await setup({
    defaultPage: 'inventory',
    baseScriptPath: '/static/js',
    routerBasePath: '/',
    openIdConfig: creastinaOpenIdConfig,
    storagePrefix: '/creastina/crafting',
  });
});
