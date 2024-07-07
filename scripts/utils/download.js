import chalk from 'chalk';

export default async function download(url) {
    const res = await fetch(url);
    if (!res.ok) {
        console.log(
            chalk.bold.bgRed(' ERROR: '),
            chalk.bold.redBright(res.statusText),
        );
        process.exit(1)
    }
    return Buffer.from(await res.arrayBuffer());
}
