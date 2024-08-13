import { mkdir } from 'node:fs/promises';
import { error } from './logs.js';

/**
 * @param { string } path
 */

export default async function create(path) {
    await mkdir(path, { recursive: true }).catch((err) => {
        error(err);
    });
}
