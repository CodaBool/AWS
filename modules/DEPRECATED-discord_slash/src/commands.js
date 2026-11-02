// any changes here must be registerd
// you can run a register with `npm run register`

export default {
  ADD_TO_ALLOWLIST_COMMAND: {
    // 'name' must not contain spaces
    name: 'minecraft_allowlist',
    description: 'add a minecraft username to the allowlist',
    options: [{
      name: "username",
      type: 3, // STRING
      description: "add a minecraft username to the allowlist",
      required: true,
    }],
  },
  FACTORIO_ADD_TO_ALLOWLIST_COMMAND: {
    // 'name' must not contain spaces
    name: 'factorio_allowlist',
    description: 'add a factorio username to the allowlist',
    options: [{
      name: "input",
      type: 3, // STRING
      description: "add a factorio username to the allowlist",
      required: true,
    }],
  },
  FACTORIO_ADMIN_COMMAND: {
    // 'name' must not contain spaces
    name: 'factorio_command',
    description: 'run a factorio console command',
    options: [{
      name: "input",
      type: 3, // STRING
      description: "run a factorio console command",
      required: true,
    }],
  },
}
