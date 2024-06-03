import * as dotenv from 'dotenv'
import { verifyKey, InteractionType, InteractionResponseType } from 'discord-interactions'
import { Client } from 'ssh2'
import CMD from "./commands.js"

// global Discord username
const ADMIN = "CodaBool"

dotenv.config()

export async function handler(event, context) {
  const signature = event.headers['x-signature-ed25519']
  const timestamp = event.headers['x-signature-timestamp']
  const text = event.body
  const valid = verifyKey(text, signature, timestamp, process.env.PUBLIC_KEY)

  if (!valid) {
    return { statusCode: 403, body: "unauthorized" }
  }

  const body = JSON.parse(text)
  
  if (body.type === InteractionType.PING) {
    return { statusCode: 200, body: JSON.stringify({ type: InteractionResponseType.PONG }) }
  }

  if (body.type !== InteractionType.APPLICATION_COMMAND) {
    return { statusCode: 502, body: JSON.stringify({ error: 'rejecting non-command interaction' }) }
  }

  console.log(body)

  if (body.data.name === CMD.ADD_TO_ALLOWLIST_COMMAND.name) {

    // only allow admin
    // if (ADMIN !== body.member.user.global_name) {
    //   return {
    //     statusCode: 200,
    //     body: JSON.stringify({
    //       type: InteractionResponseType.CHANNEL_MESSAGE_WITH_SOURCE,
    //       data: { content: `only ${ADMIN} is able to use this command` },
    //     })
    //   }
    // }

    // sanitize input to prevent shell injection
    const value = body.data.options.find(o => o.name === "username").value
    const sanitizedValue = value.replace(/[^a-zA-Z0-9_.-]/g, '')
    if (value !== sanitizedValue) {
      return {
        statusCode: 200,
        body: JSON.stringify({
          type: InteractionResponseType.CHANNEL_MESSAGE_WITH_SOURCE,
          data: { content: `your username contains illegal characters ${value}` },
        })
      };
    }

    // SSH
    const content = await new Promise(function (resolve, reject) {
      const conn = new Client()

      conn.on('ready', () => {
        conn.exec(`mcrcon -p ${process.env.RCON_PASSWORD} "whitelist add ${sanitizedValue}"`, (err, stream) => {
          if (err) throw err
          stream.on('close', conn.end)
            .on('data', stdout => {
              console.log("raw output =", stdout)
              console.log(stdout.toString().replace(/\x1b\[[0-9;]*m/g, ''))
              resolve(stdout.toString().replace(/\x1b\[[0-9;]*m/g, ''))
            }).stderr.on('data', console.log)
        })
      }).connect({
        host: process.env.SERVER_IP,
        port: process.env.SSH_PORT,
        username: process.env.SSH_USERNAME,
        password: process.env.SSH_PASSWORD
      })
    })

    // respond back with stdout
    return {
      statusCode: 200,
      body: JSON.stringify({
        type: InteractionResponseType.CHANNEL_MESSAGE_WITH_SOURCE,
        data: { content },
      })
    }
  }

  return { statusCode: 404, body: "not found" };
}
