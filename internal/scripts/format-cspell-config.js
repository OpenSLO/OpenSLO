/*
  Formatter which works on cspell config file and:
  - Sorts the 'words' list.
  - Removes duplicates from 'words' list.
*/

import { readFileSync, writeFileSync } from 'fs';

const CSPELL_CONFIG = "cspell.json"

function format() {
  const f = readFileSync(CSPELL_CONFIG, 'utf8')
  const contents = JSON.parse(f, { keepSourceTokens: true })

  let words = contents['words']
  let set = new Set()
  words = words.sort().filter((word) => {
    if (!set.has(word)) {
      set.add(word)
      return true
    }
    return false
  })
  contents['words'] = words

  writeFileSync(CSPELL_CONFIG, JSON.stringify(contents, null, 2))
}

try { format() } catch (err) { console.error(err) }
