import { mkdir } from 'fs/promises';

export default async function create(path) {
    await mkdir(path, { recursive: true }).catch((err) => {
        console.log(chalk.bold.bgRed(' ERROR: '), chalk.bold.redBright(err));
        process.exit(1);
    });
}
