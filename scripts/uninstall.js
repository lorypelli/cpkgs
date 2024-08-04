import { existsSync } from 'fs';
import { file } from './utils/constants.js';
import { info, success, warning } from './utils/logs.js';
import remove from './utils/remove.js';

if (existsSync(file)) {
    info(`Removing ${file}...`);
    await remove(file);
    success('File successfully removed!');
} else {
    warning('Directory does not exists, nothing was changed!');
}
