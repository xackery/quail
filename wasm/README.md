# Quail WASM

Quail WASM is a port of Quail to WebAssembly and meant to be used in the browser (node.js support not considered).

There is an example under examples/index.html and that can be run via `npm install` and `npm start` from this folder

The entrypoint for using this module is as follows

```js
import { CreateQuail } from 'quail-wasm'

// Depending on where you are hosting quail.wasm
CreateQuail('/static/quail.wasm').then(q => {
    const { quail, fs } = q;

    // Convert from s3d to json
    fs.setEntry('/qeynos2.s3d', /* Uint8Array */ fs.makeFileEntry(undefined, someBuffer));
    quail.convert('/qeynos2.s3d', '/qeynos2.json');

    const qeynos2 = JSON.parse(new TextDecoder().decode(fs.getEntry('/qeynos2.json').data, 'utf8'));

    // Write from json back to s3d
    quail.convert('/qeynos2.json', '/qeynos2_new.s3d');
});
```