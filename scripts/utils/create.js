import chalk from 'chalk';
import { mkdir } from 'fs';

export default function create(path) {
    mkdir(path, { recursive: true }, (err) => {
        if (err) {
            console.log(chalk.bold.bgRed('ERROR:'), chalk.bold.redBright(err));
        }
    });
}