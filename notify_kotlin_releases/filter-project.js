const { resolve } = require('path')
const { writeFileSync } = require('fs')
const file = resolve(__dirname, 'current.json')
const releases = require(file)
writeFileSync(file, JSON.stringify(
  releases
    .filter(({ draft }) => !draft)
    .map(({
      id,
      name,
      html_url,
      prerelease,
      assets,
    }) => ({
      id,
      name,
      html_url,
      prerelease,
      download_url: assets && assets[0] && assets[0].name.startsWith('kotlin-compiler-') ? assets[0].browser_download_url : null,
    }))
))
