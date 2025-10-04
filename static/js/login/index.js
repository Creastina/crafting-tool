import { Alpine } from '../../lib/alpine.js';
import { openIdLogin } from '../../lib/jinya-alpine-tools.js';

Alpine.data('loginData', () => ({
  login: openIdLogin,
}));
