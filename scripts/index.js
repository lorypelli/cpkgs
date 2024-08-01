import { existsSync } from 'fs';
import create from './utils/create.js';
import download from './utils/download.js';
import { info, success, warning } from './utils/logs.js';
import write from './utils/write.js';

const extension = process.platform == 'win32' ? '.exe' : '';
const url = `https://github.com/lorypelli/cpkgs/releases/latest/download/cpkgs_${process.platform}${extension}`;
info('Downloading...');
const buffer = await download(url);
success('File downloaded successfully!');
const dir =
    process.platform == 'win32'
        ? `${process.env.APPDATA}/cpkgs/bin`
        : '/usr/local/bin/cpkgs';
if (!existsSync(dir)) {
    warning('Directory does not exists!');
    info('Creating...');
    await create(dir);
    success('Directory created successfully!');
}
const file = `${dir}/cpkgs${extension}`;
info('Writing...');
await write(file, buffer);
success('File written successfully!');
success('cpkgs has been installed correctly');
warning(
    `please make sure that ${dir} is under system PATH and restart your shell!`,
);
