let knownHashes = {}

function getLatestHash(path) {
    let request = new XMLHttpRequest()
    request.open("GET", `/hash/md5${path}`)

    request.onload = () => {
        if(request.status === 200) {
            let hashResp = JSON.parse(request.responseText)
            if(path in knownHashes && knownHashes[path] !== hashResp["Hash"]) {
                window.location.reload()
            } else {
                knownHashes[path] = hashResp["Hash"]
            }
        }
    }
    request.send()
}

function subscribeForUpdates(path) {
    setInterval(() => {
        getLatestHash(path)
    }, 1000)
}