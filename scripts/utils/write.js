import chalk from 'chalk';
import { writeFile } from 'fs';

export default function write(path, buffer) {
    writeFile(path, buffer, (err) => {
        if (err) {
            console.log(chalk.bold.bgRed('ERROR:'), chalk.bold.redBright(err));
        }
    });
}
