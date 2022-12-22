import jsdom from 'jsdom'
import axios from 'axios'
import pg from 'pg'
import format from 'pg-format'
import * as dotenv from 'dotenv'
import pypi from 'pypi-info'
dotenv.config()

const { JSDOM } = jsdom

const DOCKER = process.env.AWS_LAMBDA_FUNCTION_NAME

export async function handler(event, context) {
  const cfg = {}
  process.argv.indexOf('-t') > -1 ? cfg['test'] = true : cfg['test'] = false

  const argIndex = process.argv.indexOf('--path')
  if (argIndex > -1) {
    console.log('set path to', process.argv[argIndex + 1])
    event['pathParameters']['id'] = process.argv[argIndex + 1]
  }

  if (DOCKER) {
    if (DOCKER == 'test_function') { 
      // local docker
      event = {
        path: "/v1/upcoming_movies",
        queryStringParameters: {
          limit: 25
        }
      }
    } else { 
      // cloud docker
    }
  }

  let response = { 
    statusCode: 200, 
    body: 'default',
    headers: {
      "Access-Control-Allow-Headers": "*",
      "Access-Control-Allow-Origin": "*",
      "Access-Control-Allow-Methods": "*"
    },
  }
  const db = new pg.Client({ 
    connectionString: process.env.PG_URI,
    ssl: { rejectUnauthorized: false }
  })
  try {    
    let creationKey = ''
    let limit = 100

    if (event["Records"] !== undefined) {
      event.body = event["Records"][0].body
      // console.log('body', event.body)
      event.path = event.body.split('@')[0]
      creationKey = event.body.split('@')[1]
      if (creationKey === process.env.KEY) {
        console.log('Triggered from SQS with an authorized write request')
      }
    } else {
      console.log('Trigger from outside SQS')
      creationKey = event.queryStringParameters?.key
      limit = event.queryStringParameters?.limit
    }

    console.log('query', event.queryStringParameters)
    console.log('path', event.path)
    console.log('limit', limit)
    // console.log('body', event.body)

    if (creationKey) { // write
      if (creationKey !== process.env.KEY) throw 'Wrong key'

      if (event.path === '/v1/trending_github') {
        if (!process.env.GIT_TOKEN) throw 'undefined GIT_TOKEN env var'
        response.body = await githubTrends()
      } else if (event.path === '/v1/upcoming_movies') {
        response.body = await getUpComingMovies()
      } else if (event.path === '/v1/trending_movies') {
        response.body = await getTrendingMovie()
      } else if (event.path === '/v1/trending_tv') {
        response.body = await getTrendingTV()
      } else if (event.path === '/v1/trending_games') {
        response.body = await trendingGames()
      } else if (event.path === '/v1/trending_npm') {
        response.body = await getNpmTrend()
      } else if (event.path === '/v1/trending_pypi') {
        response.body = await getPyTrend()
      } else if (event.path === '/v1/get_build') {
        response.body = process.env.BUILD_ID
      } else {
        throw `BUILD: ${process.env.BUILD_ID} |
  Use one of the following api paths:
  /trending_github
  /upcoming_movies
  /trending_movies
  /trending_tv
  /trending_games
  /trending_npm
  /trending_pypi`
      }
  
      if (response.body && (event.path !== '/v1/get_build')) { // write
        console.log('save to db')
        await db.connect()
        let { deleteSQL, insertSQL } = generateSQL(event.path, response.body)
        console.log('SQL DUMP', deleteSQL)
        await db.query(deleteSQL)
        console.log('insert sql', insertSQL)
        const res = await db.query(insertSQL).then(res => res.rowCount)
        console.log('db query res', res)
      }
    } else { // read
      console.log('read request')
      await db.connect()
      const readSQL = generateSQL(event.path, null, true, limit)
      console.log('read sql =', readSQL)
      if (!readSQL) throw `BUILD: ${process.env.BUILD_ID} |
Use one of the following api paths:
/trending_github
/upcoming_movies
/trending_movies
/trending_tv
/trending_games
/trending_npm
/trending_pypi`
      response.body = await db.query(readSQL).then(res => res.rows)
      console.log('rows', response.body)
    }
    response.body = JSON.stringify(response.body, null, 2)
  } catch (err) {
    console.log(err)
    if (typeof err === 'string') {
      response = { statusCode: 400, body: err.split('\n').join(' ') }
    } else {
      response = { statusCode: 500, body: (err.message || err)}
    }
  } finally { 
    await db.end()
    return response
  }
}

if (!DOCKER) handler({
  path: "/v1/upcoming_movies",
  queryStringParameters: {
    limit: 25
  }
})

function toArr(rawData) {
  return rawData.map(obj => {
    // console.log('------------------')
    return Object.keys(obj).map(key => { 
      // console.log(key)
      return obj[key]
    })
  })
}

function generateSQL(path, data, isRead, limit) {
  let deleteSQL = null
  let insertSQL = null
  let readSQL = null
  // TODO: limit should be set to high number and given to the format method for read
  if (path === '/v1/trending_github') {
    deleteSQL = 'DELETE FROM trending_github'
    insertSQL = 'trending_github(name, href, description, stars)'
    readSQL = `SELECT * FROM trending_github${limit ? ` LIMIT ${limit}` : ''}`
  } else if (path === '/v1/trending_npm') {
    deleteSQL = 'DELETE FROM trending_npm'
    insertSQL = 'trending_npm(subject, page, rank, title, description)'
    readSQL = 'SELECT * FROM trending_npm ORDER BY subject, rank'
  } else if (path === '/v1/trending_movies') {
    deleteSQL = 'DELETE FROM trending_movies'
    insertSQL = 'trending_movies(link, img, title, year, rank, velocity, rating)'
    readSQL = `SELECT * FROM trending_movies${limit ? ` LIMIT ${limit}` : ''}`
  } else if (path === '/v1/trending_pypi') {
    deleteSQL = 'DELETE FROM trending_pypi'
    insertSQL = 'trending_pypi(name, description, downloads)'
    readSQL = `SELECT * FROM trending_pypi${limit ? ` LIMIT ${limit}` : ''}`
  } else if (path === '/v1/trending_tv') {
    deleteSQL = 'DELETE FROM trending_tv'
    insertSQL = 'trending_tv(link, img, title, rank, velocity, rating)'
    readSQL = `SELECT * FROM trending_tv${limit ? ` LIMIT ${limit}` : ''}`
  } else if (path === '/v1/trending_games') {
    deleteSQL = 'DELETE FROM trending_games'
    insertSQL = 'trending_games(link, title, rating, is_free, regular_price, discounted_price)'
    readSQL = `SELECT * FROM trending_games${limit ? ` LIMIT ${limit}` : ''}`
  } else if (path === '/v1/upcoming_movies') {
    deleteSQL = 'DELETE FROM upcoming_movies'
    insertSQL = 'upcoming_movies(title, release)'
    readSQL = `SELECT * FROM upcoming_movies${limit ? ` LIMIT ${limit}` : ''}`
    // for mongo like object data I can use jsonb,
    // format(`INSERT INTO table(raw_json) VALUES(%L)`, [JSON.stringify(data)])
    // if (isRead) return `SELECT * FROM upcoming_movies${limit ? ` LIMIT ${limit}` : ''}`
    // return {
    //   deleteSQL: 'DELETE FROM upcoming_movies',
    //   insertSQL: format(`INSERT INTO upcoming_movies(raw_json) VALUES(%L)`, [JSON.stringify(data)])
    // }
  } else return
  if (isRead) return readSQL
  return {
    deleteSQL,
    insertSQL: format(`INSERT INTO ${insertSQL} VALUES %L`, toArr(data))
  }
}

async function githubTrends() {
  const LANGUAGES = ["JavaScript", "Python", "Shell"]
  const TOKEN = process.env.GIT_TOKEN;
  const allData = await axios.get('https://api.github.com/search/repositories?q=stars:>0&sort=stars&per_page=100', 
      { Headers: { Authorization: TOKEN } })
    .then(res => res.data)
    .catch(console.log)
  for (const language of LANGUAGES) {
    const langData = await axios.get(`https://api.github.com/search/repositories?q=language:${language}&stars:>0&sort=stars&per_page=100`, 
        { Headers: { Authorization: TOKEN } })
      .then(res => res.data)
      .catch(console.log)
    const relevantLangData = langData.items.map(repo => ({
      name: repo.name,
      href: repo.html_url,
      description: repo.description?.substring(0, 300) || '',
      stars: repo.stargazers_count,
      language
    }))
    console.log(relevantLangData)
  }
  const relevantAllData = allData.items.map(repo => ({
    name: repo.name,
    href: repo.html_url,
    description: repo.description.substring(0, 300),
    stars: repo.stargazers_count
  }))
  return relevantAllData
}

async function getUpComingMovies() {
  // IMDB will give 403 to prevent scraping unless a User-Agent is spoofed
  // seems to assume us and movie
  const html = await axios.get('https://www.imdb.com/calendar/?region=US&type=MOVIE', {
    headers: { 'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:108.0) Gecko/20100101 Firefox/108.0' }
  })
    .then(res => res.data)
    .catch(err => console.log('bad request', err))
  const dom = new JSDOM(html)
  const articles = dom.window.document.querySelector(".ipc-page-section--base")?.getElementsByTagName('article')
  const data = []
  for (const a of articles) {
    const release = a?.firstChild?.textContent.trim()
    console.log('==> ', release)
    for (const li of a?.lastChild.childNodes) {
      const title = li?.textContent.trim().split('(20')[0]
      data.push({title, release})
      console.log('           ', title)
    }
  }
  return data
}

async function getTrendingTV() {
  const html = await axios.get('https://www.imdb.com/chart/tvmeter')
    .then(res => res.data)
    .catch(console.log)
  const dom = new JSDOM(html)
  const list = dom.window.document.querySelector(".lister-list")
  const data = []
  for (const row of list.getElementsByTagName('tr')) {
    const show = {}
    for (const node of row.childNodes) {
      if (node.className === 'posterColumn') {
        show.link = 'https://imdb.com' + node.getElementsByTagName('a').item(0).getAttribute('href')
        show.img = node.getElementsByTagName('img').item(0).getAttribute('src')
      } else if (node.className === 'titleColumn') {
        show.title = node.getElementsByTagName('a').item(0).textContent
        const el = node.getElementsByTagName('div').item(0).querySelector('.secondaryInfo')?.textContent
        let velocity = 0
        if (el) {
          velocity = Number(el.split('\n')[2].slice(0, -1).match(/\d+/g).join(''))
        }
        let rank = node.getElementsByTagName('div').item(0).childNodes[0].textContent.trim()
        if (rank.includes('no change')) {
          rank = rank.match(/\d+/g).toString() // remove (no change) text
          velocity = 0
        }

        if (node.getElementsByTagName('div').item(0).childNodes[1]?.childNodes[1]?.className.includes('down')) {
          velocity *= -1 // use negative to represent downward velocity
        }
        show.rank = rank
        show.velocity = velocity
      } else if (node.className?.includes('imdbRating')) {
        show.rating = node.textContent.trim()
      }
    }
    console.log('+', show)
    data.push(show)
  }
  return data
}

async function getTrendingMovie() {
  const html = await axios.get('https://www.imdb.com/chart/moviemeter')
    .then(res => res.data)
    .catch(console.log)
  const dom = new JSDOM(html)
  const list = dom.window.document.querySelector(".lister-list")
  const data = []
  for (const row of list.getElementsByTagName('tr')) {
    const movie = {}
    for (const node of row.childNodes) {
      if (node.className === 'posterColumn') {
        movie.link = 'https://imdb.com' + node.getElementsByTagName('a').item(0).getAttribute('href')
        movie.img = node.getElementsByTagName('img').item(0).getAttribute('src')
      } else if (node.className === 'titleColumn') {
        movie.title = node.getElementsByTagName('a').item(0).textContent
        movie.year = node.getElementsByTagName('span').item(0).textContent.split('').slice(1, 5).join('')
        const el = node.getElementsByTagName('div').item(0).querySelector('.secondaryInfo')?.textContent
        let velocity = 0
        if (el) {
          velocity = Number(el.split('\n')[2].slice(0, -1).match(/\d+/g).join(''))
        }
        let rank = node.getElementsByTagName('div').item(0).childNodes[0].textContent.trim()
        if (rank.includes('no change')) {
          rank = rank.match(/\d+/g).toString() // remove (no change) text
          velocity = 0
        }

        if (node.getElementsByTagName('div').item(0).childNodes[1]?.childNodes[1]?.className.includes('down')) {
          velocity *= -1 // use negative to represent downward velocity
        }
        movie.rank = rank
        movie.velocity = velocity
      } else if (node.className?.includes('imdbRating')) {
        movie.rating = node.textContent.trim()
      }
    }
    console.log('+', movie)
    data.push(movie)
  }
  return data
}

async function trendingGames() {
  const html = await axios.get('https://store.steampowered.com/search/?filter=topsellers')
    .then(res => res.data)
    .catch(console.log)
  const dom = new JSDOM(html)
  const list = dom.window.document.querySelector("#search_resultsRows")
  const data = []
  for (const row of list.getElementsByTagName('a')) {
    const game = {}
    game.link = null
    game.title = null
    game.rating = null
    game.is_free = false
    game.regular_price = null
    game.discounted_price = null
    game.link = row.getAttribute('href')
    const ratingRaw = row.childNodes[3].querySelector('.search_review_summary')?.getAttribute('data-tooltip-html')
    if (ratingRaw) game.rating = ratingRaw.split('%')[0].split('<br>')[1] + '%'
    game.title = row.childNodes[3].querySelector('.title').textContent.trim()
    const arr = row.childNodes[3].querySelector('.search_price').textContent.trim().split('$')
    if (arr.length === 3) {
      game.discounted_price = Number(arr[2].replace(/[^0-9.-]+/g,""))
      game.regular_price = Number(arr[1].replace(/[^0-9.-]+/g,""))
    } else if (arr.length === 2) {
      game.regular_price = Number(arr[1].replace(/[^0-9.-]+/g,""))
    } else {
      game.is_free = true
    }
    data.push(game)
  }
  return data
}

async function getNpmTrend() {
  const keywords = [
    'backend',
    'front-end',
    'cli',
    'framework',
  ]
  const data = []
  for (let page = 0; page < 2; page++) {
    for (const subject of keywords) {
      let count = 0
      console.log('\n' +subject, '| page', page)
      console.log('scraping', `https://www.npmjs.com/search?ranking=popularity&page=${page}&q=keywords%3A${subject}`)
      const html = await axios.get(`https://www.npmjs.com/search?ranking=popularity&page=${page}&q=keywords%3A${subject}`, { 
        headers: { "Accept-Encoding": "gzip,deflate,compress" } 
      }).then(res => res.data)
        .catch(console.log)
      const dom = new JSDOM(html)
      const list = dom.window.document.getElementsByTagName("main").item(0)?.childNodes[2]?.childNodes[1]
      if (!list) {
        console.log('no results found')
        return
      }
      for (const result of list.childNodes) { // section

        if (result.firstChild.firstChild.firstChild.childNodes.length === 2) {
          let rank = (page * 20) + count + 1
          const item = { subject, page, rank }
          item.title = result.firstChild.firstChild.firstChild.lastChild.textContent
          // console.log('+', item.title)
          item.description = 'None' // some npm have no descrition eg. @devoralime/server
          for (const node of result.firstChild.childNodes) { // all nodes inside 1st div of results

            // if (item.title === '@devoralime/server') {
            //   console.log('raw node', node)
            //   console.log('node text', node.textContent)
            //   console.log('node name', node.nodeName)
            // }
            if (node.nodeName === 'P') { // description
              item.description = node.textContent
            }
          }
          // console.log(`page ${page} * 10 + count ${count} + 1 = `, rank)
          count++
          data.push(item)
        }
      }
      console.log('finished', subject, 'page', page, 'with', data.length, 'items')
    }
  }
  // console.log(data)
  return data
}

async function getPyTrend() {
  const data = {}
  console.log('scraping python packages')
  const json = await axios.get('https://hugovk.github.io/top-pypi-packages/top-pypi-packages-30-days.min.json')
    .then(res => res.data)
    .catch(err => console.log(err))
  console.log('found the top', json.rows.length, 'most popular Python packages')
  for (let i = 0; i < 100; i ++) {
    data[json.rows[i].project] = {downloads: json.rows[i].download_count, description: '🤷'}
  }
  const newPromises = []
  const issues = []
  for (const [key, value] of Object.entries(data)) {
    newPromises.push(new Promise((resolve, reject) => {
      pypi.getPackage(key)
        .then(res => {
          resolve({name: key, description: res?.info?.summary})
          // I doubt I would need to return here
          return {name: key, description: res?.info?.summary}
        })
        .catch(err => {
          console.log('err for', key)
          issues.push(key)
          reject(key)
        })
    }))
  }
  
  console.log('I am limiting this to the top', Object.keys(data).length, 'requesting a summary for each now...')
  await Promise.allSettled(newPromises)
    .then(results => results.forEach(result => {
      if (result?.status == 'fulfilled') {
        data[result?.value?.name]['description'] = result?.value?.description
      }
    }))
    .finally(() => console.log('Failed getting descriptions for', issues.length))
  
  // flatten data
  const flat_data = []
  for (const [key, value] of Object.entries(data)) {
    flat_data.push({
      name: key,
      description: value.description,
      downloads: value.downloads
    })
  }
  return flat_data
}