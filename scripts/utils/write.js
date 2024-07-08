import chalk from 'chalk';
import { writeFile } from 'fs/promises';

/**
 * @param { string } path
 * @param { string } buffer
 */

export default async function write(path, buffer) {
    await writeFile(path, buffer).catch((err) => {
        console.log(chalk.bold.bgRed(' ERROR: '), chalk.bold.redBright(err));
        process.exit(1);
    });
}
