import chalk from 'chalk';
import { existsSync } from 'fs';
import create from './utils/create.js';
import download from './utils/download.js';
import write from './utils/write.js';

const platform = process.platform == 'win32' ? 'windows' : process.platform;
const extension = platform == 'windows' ? '.exe' : '';
const url = `https://github.com/lorypelli/cpkgs/releases/latest/download/cpkgs_${platform}${extension}`;
console.log(chalk.bold.blueBright('Downloading...'));
const buffer = await download(url);
console.log(
    chalk.bold.bgGreen(' SUCCESS: '),
    chalk.bold.greenBright('File downloaded successfully!'),
);
const dir =
    platform == 'windows'
        ? `${process.env.APPDATA}/cpkgs`
        : '/usr/local/bin/cpkgs';
if (!existsSync(dir)) {
    console.log(
        chalk.bold.bgYellow(' WARNING: '),
        chalk.bold.yellowBright('Directory does not exists!'),
    );
    console.log(chalk.bold.blueBright('Creating...'));
    await create(dir);
    console.log(
        chalk.bold.bgGreen(' SUCCESS: '),
        chalk.bold.greenBright('Directory created successfully!'),
    );
}
const file = `${dir}/${url.split('/').at(-1)}`;
console.log(chalk.bold.blueBright('Writing...'));
await write(file, buffer);
console.log(
    chalk.bold.bgGreen(' SUCCESS: '),
    chalk.bold.greenBright('File written successfully!'),
);
