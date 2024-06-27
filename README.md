# CPKGS

A CLI application to easily install C modules and run your code fast.

## Usage

You have multiple commands avaible. To see them all type: `cpkgs help`

### Initializing a new project

To do this you would need to use the `cpkgs init` command, you will be asked about which compiler to use and how the output filename should be called

### Installing a package

To do this you would need to use the `cpkgs add` command. You need to provide github repo of the package and header file name.

### Running the project

Finally, when you are all done, to run the project you need to use the `cpkgs run` command with the file you want to run passed as an argument.

### Removing the project

You aren't satisfied by it, no problem, you can use the `cpkgs remove` command by providing directory name and it will be completely removed.

## Installing the entire project

By using the `cpkgs install` command, all of the packages will be fetched from the `cpkgs.json` file and will be installed automatically.

<b>Note:</b> the `cpkgs` directory <b>SHOULD BE ALWAYS</b> placed in the `.gitignore` file.

I don't do this on the `example` directory on purpose only to showcase how the entire project will look but you <b>SHOULD ALWAYS</b> do this.

## Updating packages

By using the `cpkgs update` command you can update a specific package or by providing `-a` flag you can update all packages at once.