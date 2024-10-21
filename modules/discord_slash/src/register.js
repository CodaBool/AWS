import COMMANDS from './commands.js'
import dotenv from 'dotenv'

/**
 * This file is meant to be run from the command line, and is not used by the
 * application server. Only needs to be run once.
 */

dotenv.config()

if (!process.env.TOKEN) {
  throw new Error('The TOKEN environment variable is required.')
}
if (!process.env.APP_ID) {
  throw new Error('The APP_ID environment variable is required.')
}

/**
 * Register all commands globally. This can take o(minutes), so wait until
 * you're sure these are the commands you want.
 *
 * https://discord.com/developers/docs/interactions/application-commands
 */
const url = `https://discord.com/api/v10/applications/${process.env.APP_ID}/commands`

console.log("registering commands", Object.values(COMMANDS).map(cmd => cmd.name))

const response = await fetch(url, {
  headers: {
    'Content-Type': 'application/json',
    Authorization: `Bot ${process.env.TOKEN}`,
  },
  method: 'PUT',
  body: JSON.stringify(Object.values(COMMANDS)),
});

if (response.ok) {
  console.log('Registered all commands');
  const data = await response.json();
  console.log(JSON.stringify(data, null, 2));
} else {
  console.error('Error registering commands');
  let errorText = `Error registering commands \n ${response.url}: ${response.status} ${response.statusText}`;
  try {
    const error = await response.text();
    if (error) {
      errorText = `${errorText} \n\n ${error}`;
    }
  } catch (err) {
    console.error('Error reading body from request:', err);
  }
  console.error(errorText);
}
