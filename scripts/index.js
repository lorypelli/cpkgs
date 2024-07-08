import chalk from 'chalk';
import { existsSync } from 'fs';
import create from './utils/create.js';
import download from './utils/download.js';
import write from './utils/write.js';

const extension = process.platform == 'win32' ? '.exe' : '';
const url = `https://github.com/lorypelli/cpkgs/releases/latest/download/cpkgs_${process.platform}${extension}`;
console.log(chalk.bold.blueBright('Downloading...'));
const buffer = await download(url);
console.log(
    chalk.bold.bgGreen(' SUCCESS: '),
    chalk.bold.greenBright('File downloaded successfully!'),
);
const dir =
    process.platform == 'win32'
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
const file = `${dir}/cpkgs${extension}`;
console.log(chalk.bold.blueBright('Writing...'));
await write(file, buffer);
console.log(
    chalk.bold.bgGreen(' SUCCESS: '),
    chalk.bold.greenBright('File written successfully!'),
);
console.log(
    chalk.bold.bgYellow(' WARNING: '),
    chalk.bold.yellowBright(
        `cpkgs has been installed correctly, if you are not able to use it, please make sure that ${dir} is under system PATH and restart your shell!`,
    ),
);
