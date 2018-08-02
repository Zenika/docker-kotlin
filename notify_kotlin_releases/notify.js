const { request } = require('https')
const { resolve } = require('path')
const previous = require(resolve(__dirname, 'previous.json'))
const current = require(resolve(__dirname, 'current.json'))
const previousIds = previous.map(({ id }) => id)
current
  .filter(({ id }) => !previousIds.includes(id))
  .forEach(({ name, html_url, prerelease, download_url }) => {
    request({
      host: 'api.github.com',
      method: 'POST',
      path: '/repos/Zenika/alpine-kotlin/issues',
      auth: `${process.env.GITHUB_ACCOUNT}:${process.env.GITHUB_OAUTH_TOKEN}`,
      headers: {
        'User-Agent': 'CircleCI'
      }
    }, res => {
      if (res.statusCode >= 200 && res.statusCode <= 299) return
      const error = ''
      res.on('data', chunk => {
        error += chunk
      })
      res.on('end', () => console.error(error))
      process.exit(1)
    })
      .end(JSON.stringify({
        title: `⬆️ New Kotlin release ${name}`,
        body: `[Kotlin version ${name}](${html_url}) has been released.${prerelease ? '\n⚠ This is a pre-release.' : ''}\nCompiler: ${download_url}`,
        labels: ['Kotlin release']
      }))
  })
