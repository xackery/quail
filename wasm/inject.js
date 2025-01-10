const fs = require('fs');
const package = require('./package.json');
package.version = process.env.VERSION;
fs.writeFileSync('./package.json', JSON.stringify(package, null, 4));
console.log(`Injected version into package.json ${package.version}`);
