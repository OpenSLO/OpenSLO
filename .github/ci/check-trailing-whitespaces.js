/*
    Linter which checks if all files under git control do not contain any trailing
    white spaces (both spaces and tabs characters), moreover non-text files are
    excluded from check based on extension from array fileExtensionsToIgnore
    Requires git available in PATH and can be run only in a repository
*/

const fs = require('fs')
const { spawnSync } = require('child_process')

const fileExtensionsToIgnore = ['.ico', '.png', '.desc']

// get all files under git control
const gitListFiles = spawnSync('git', ['ls-tree', '-r', 'HEAD', '--name-only'])
if (gitListFiles.stderr.toString() !== "") {
    console.error(`Unexpected error occurred: ${gitListFiles.stderr.toString()}`)
    process.exit(2)
}

const filesToCheck = gitListFiles.stdout.toString().split('\n').filter(file =>
    fileExtensionsToIgnore.every(extension => !file.endsWith(extension))
)

const noTrailingWhitespaces = new RegExp(/[ \t]+$/gm)
filesToCheck.forEach(file => {
    fs.readFile(file, 'utf8', (_, content) => {
        const match = noTrailingWhitespaces.exec(content);
        if (match) {
            console.error(`${file} contains trailing whitespaces around: `,
                match.input.substr(Math.max(0, match.index - 30), 60));
            process.exitCode = 1
        }
    })
})
