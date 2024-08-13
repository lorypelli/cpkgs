import { existsSync } from 'node:fs';
import { dir, extension, file } from './utils/constants.js';
import create from './utils/create.js';
import download from './utils/download.js';
import { info, success, warning } from './utils/logs.js';
import write from './utils/write.js';

const url = `https://github.com/lorypelli/cpkgs/releases/latest/download/cpkgs_${process.platform}${extension}`;
info('Downloading...');
const buffer = await download(url);
success('File downloaded successfully!');
if (!existsSync(dir)) {
    warning('Directory does not exists!');
    info('Creating...');
    await create(dir);
    success('Directory created successfully!');
}
info('Writing...');
await write(file, buffer);
success('File written successfully!');
success('cpkgs has been installed correctly');
warning(
    `please make sure that ${dir} is under system PATH and restart your shell!`,
);
if (process.platform != 'win32') {
    warning(
        `You may also need to give executable permissions to ${file}, use the command 'chmod +x ${file}'`,
    );
}
