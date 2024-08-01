import chalk from 'chalk';

/**
 * @param { string } msg
 */

export function error(msg) {
    console.log(chalk.bold.bgRed('  ERROR  '), chalk.bold.redBright(msg));
    process.exit(1);
}

/**
 * @param { string } msg
 */

export function info(msg) {
    console.log(chalk.bold.bgBlue('  INFO  '), chalk.bold.blueBright(msg));
}

/**
 * @param { string } msg
 */

export function success(msg) {
    console.log(chalk.bold.bgGreen('  SUCCESS  '), chalk.bold.greenBright(msg));
}

/**
 * @param { string } msg
 */

export function warning(msg) {
    console.log(
        chalk.bold.bgYellow('  WARNING  '),
        chalk.bold.yellowBright(msg),
    );
}
