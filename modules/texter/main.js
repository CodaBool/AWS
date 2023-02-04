import { SNSClient, PublishCommand } from "@aws-sdk/client-sns"
import * as dotenv from 'dotenv'
import fetch from 'node-fetch'
import twilio from 'twilio'
// import fs from 'fs'

dotenv.config()
const client = new twilio(process.env.TWILIO_ACCOUNT_SID, process.env.TWILIO_AUTH_TOKEN)

const DOCKER = process.env.AWS_LAMBDA_FUNCTION_NAME
const TWILIO_PHONE_NUMBER = process.env.TWILIO_PHONE_NUMBER
const MY_NUMBER = process.env.MY_PHONE_NUMBER

async function checkIP() {
  const lastIP = fs.readFileSync('/home/codabool/node-scripts/ip.txt', 'utf8')
  const currentIP = await fetch('http://ifconfig.me/all.json')
    .then(res => res.json())
    .then(data => data.ip_addr)
  console.log('lastIP =', lastIP)
  console.log('currentIP =', currentIP)
  if (lastIP.trim() === currentIP) {
    console.log('up to date')
    return
  }
  console.log('ip has changed')
  fs.writeFile('/home/codabool/node-scripts/ip.txt', currentIP, (err) => {
    if (err) {
      console.log(err)
      return
    }
    // fs.readFileSync('./ip.txt', 'utf8')
    console.log('ip updated')
    client.messages
      .create({
        body: 'Your Home IP has changed',
        from: TWILIO_PHONE_NUMBER,
        to: process.env.MY_PHONE_NUMBER
      })
      .catch(err => console.log(err))
  })
}

async function email(Message, Subject) {
  const sns = new SNSClient({ region: "us-east-1" })
  const command = new PublishCommand({
    Message,
    Subject,
    TopicArn: 'arn:aws:sns:us-east-1:919759177803:notify',
  })
  return await sns.send(command)
}

async function text(body) {
  return await client.messages
    .create({
      body: 'vacate',
      from: TWILIO_PHONE_NUMBER,
      to: MY_NUMBER
    })
    .then(message => console.log(message.sid))
    .done()
}

export const handler = async (event, context) => {
  if (DOCKER) {
    if (DOCKER == 'test_function') { // local docker
      console.log('local docker')
    } else { // cloud docker
      console.log('cloud docker')
    }
  }

  let data
  try {
    // data = await text("hi")
    data = await email(event.message, event.subject)
    console.log('final data', data)
  } catch (err) {
    console.log(err)
  } finally {
    return {
      statusCode: 200,
      body: JSON.stringify({ data })
    }
  }
}

if (!DOCKER) handler()