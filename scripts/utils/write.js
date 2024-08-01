import { writeFile } from 'fs/promises';
import { error } from './logs.js';

/**
 * @param { string } path
 * @param { string } buffer
 */

export default async function write(path, buffer) {
    await writeFile(path, buffer).catch((err) => {
        error(err);
    });
}
