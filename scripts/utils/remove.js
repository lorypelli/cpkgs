import { rm } from 'fs/promises';
import { error } from './logs.js';

/**
 * @param { string } path
 */

export default async function create(path) {
    await rm(path, { recursive: true }).catch((err) => {
        error(err);
    });
}
