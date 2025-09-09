window.onload = function(e) {
    const go = new Go();
    WebAssembly.instantiateStreaming(fetch('2048.wasm'), go.importObject)
        .then(result => {
            console.log(result)
            go.run(result.instance)
        })
}