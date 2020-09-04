if (!WebAssembly.instantiateStreaming) {
    // polyfill
    WebAssembly.instantiateStreaming = async (resp, importObject) => {
        const source = await (await resp).arrayBuffer();
        return await WebAssembly.instantiate(source, importObject);
    };
}

const go = new Go();
let mod, inst;
WebAssembly.instantiateStreaming(fetch("main.wasm"), go.importObject).then(
    async result => {
        mod = result.module;
        inst = result.instance;
        await go.run(inst);
    }
);

function runChip(start = true) {
    $('#dropdownMenuButton').prop('disabled', start);
    $('#runExampleBtn').prop('disabled', start);
    $('#power-indicator').toggleClass('text-success text-danger')
    if (start) {
        run();
    } else {
        stopCpu();
    }
}
