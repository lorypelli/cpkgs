import chalk from 'chalk';
import { existsSync } from 'fs';
import download from './utils/download.js';
import write from './utils/write.js';
import create from './utils/create.js';

const platform = process.platform == 'win32' ? 'windows' : process.platform;
const extension = platform == 'windows' ? '.exe' : '';
const url = `https://github.com/lorypelli/FetchTTP/releases/latest/download/FetchTTP.exe`;
console.log(chalk.bold.blueBright('Downloading...'));
const buffer = await download(url);
console.log(
    chalk.bold.bgGreen('SUCCESS:'),
    chalk.bold.greenBright('File downloaded successfully!'),
);
const dir = platform == 'windows' ? 'C:/Programs Files/cpkgs' : '/bin/cpkgs';
if (!existsSync(dir)) {
    console.log(
        chalk.bold.bgYellow('WARNING:'),
        chalk.bold.yellowBright('Directory does not exists!'),
    );
    console.log(chalk.bold.blueBright('Creating...'));
    create(dir)
    console.log(
        chalk.bold.bgGreen('SUCCESS:'),
        chalk.bold.greenBright('Directory created successfully!'),
    );
}
const file = `${dir}/${url.split('/').at(-1)}`;
console.log(chalk.bold.blueBright('Writing...'));
write(file, buffer);
console.log(
    chalk.bold.bgGreen('SUCCESS:'),
    chalk.bold.greenBright('File written successfully!'),
);
