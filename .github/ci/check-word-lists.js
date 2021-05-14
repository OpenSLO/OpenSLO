/*
    Linters which checks things related to maintaining alphabetical order of words
    in provided files, it checks the following files:
        - spelling tools dictionaries configurations:
            - cspell.json - cspell - https://github.com/streetsidesoftware/cspell/tree/master/packages/cspell
*/

const path = require('path')

const cspellConfigConfigPath = path.resolve(__dirname, '../../cspell.json')
const cspellConfig = require(cspellConfigConfigPath)
const cspellConfigWordList = {
    'list': cspellConfig.words,
    'filePath': cspellConfigConfigPath
}

const wordListsToCheck = [
    cspellConfigWordList
]
let message = ''
wordListsToCheck.forEach(wordList => {
    const loadedWordList = wordList.list
    const expectedWordList = wordList.list.slice().sort((a, b) => a.localeCompare(b))
    if (loadedWordList.length !== expectedWordList.length) {
        console.error(`Unexpected error occurred, check source code: ${__filename}`)
        process.exit(2)
    }
    for (let i = 0; i < loadedWordList.length; i++) {
        const actualWord = loadedWordList[i]
        const expectedWord = expectedWordList[i]
        if (actualWord !== expectedWord) {
            message += `${wordList.filePath}\nFirst mismatch:\n  actual: ${actualWord}\nexpected: ${expectedWord}\n\n`
            break
        }
    }
})
if (message) {
    console.error(`Alphabetical order of words is not maintained in the following files:\n\n${message}`)
    process.exit(1)
}
