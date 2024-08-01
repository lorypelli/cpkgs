import chalk from 'chalk';
import { mkdir } from 'fs/promises';

/**
 * @param { string } path
 */

export default async function create(path) {
    await mkdir(path, { recursive: true }).catch((err) => {
        console.log(chalk.bold.bgRed('  ERROR  '), chalk.bold.redBright(err));
        process.exit(1);
    });
}
