export const dir =
    process.platform == 'win32'
        ? `${process.env.APPDATA}/cpkgs/bin`
        : '/usr/bin';

export const extension = process.platform == 'win32' ? '.exe' : '';

export const file = `${dir}/cpkgs${extension}`;
