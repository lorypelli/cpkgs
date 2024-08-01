import { error } from './logs.js';

/**
 * @param { string } url
 */

export default async function download(url) {
    const res = await fetch(url);
    if (!res.ok) {
        error(res.statusText);
    }
    return Buffer.from(await res.arrayBuffer());
}
