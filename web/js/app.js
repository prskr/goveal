document.addEventListener("DOMContentLoaded", _ => {
    Promise.all([initMermaid(), setSlidesContent()])
        .then(() => {
            return initReveal()
        })
        .then(() => {
            subscribeToEvents()
        })

});

async function setSlidesContent() {
    let resp = await fetch("/slides")
    let contentText = await resp.text()
    let parser = new DOMParser()
    let contentDocument = parser.parseFromString(contentText, 'text/html')
    for (let mermaidElem of contentDocument.getElementsByClassName("mermaid")) {
        let insertSVG = (svgCode, _) => {
            mermaidElem.innerHTML = svgCode
        }

        mermaid.mermaidAPI.render('mermaid', mermaidElem.innerText, insertSVG)
    }
    document.getElementById("content-root").innerHTML = contentDocument.documentElement.innerHTML
}

async function getRevealConfig() {
    let resp = await fetch('/api/v1/config/reveal')
    return await resp.json()
}

async function initReveal() {
    let cfg = await getRevealConfig()
    Reveal.initialize({
        controls: cfg.controls,
        controlsLayout: cfg.controlsLayout,
        progress: cfg.progress,
        history: cfg.history,
        center: cfg.center,
        slideNumber: cfg.slideNumber,
        transition: cfg.transition,
        width: cfg.width,
        height: cfg.height,
        hash: true,
        pdfSeparateFragments: false,
        menu: {
            numbers: cfg.menu.numbers,
            useTextContentForMissingTitles: cfg.menu.useTextContentForMissingTitles,
            transitions: cfg.menu.transitions,
            hideMissingTitles: cfg.hideMissingTitles,
            markers: cfg.menu.markers,
            openButton: cfg.menu.openButton,
            custom: [
                {
                    title: 'Print',
                    icon: '<i class="fas fa-print"></i>',
                    content: '<a href="/?print-pdf">Go to print view<a/>'
                }
            ],
            themes: [
                {name: 'Beige', theme: '/reveal/dist/theme/beige.css'},
                {name: 'Black', theme: '/reveal/dist/theme/black.css'},
                {name: 'Blood', theme: '/reveal/dist/theme/blood.css'},
                {name: 'League', theme: '/reveal/dist/theme/league.css'},
                {name: 'Moon', theme: '/reveal/dist/theme/moon.css'},
                {name: 'Night', theme: '/reveal/dist/theme/night.css'},
                {name: 'Serif', theme: '/reveal/dist/theme/serif.css'},
                {name: 'Simple', theme: '/reveal/dist/theme/simple.css'},
                {name: 'Sky', theme: '/reveal/dist/theme/sky.css'},
                {name: 'Solarized', theme: '/reveal/dist/theme/solarized.css'},
                {name: 'White', theme: '/reveal/dist/theme/white.css'}
            ],
        },
        plugins: [RevealHighlight, RevealNotes, RevealMenu]
    })
}

async function initMermaid() {
    let resp = await fetch('/api/v1/config/mermaid')
    let cfg = await resp.json()
    mermaid.parseError = (err, hash) => {
        console.error(`Failed to parse Mermaid diagraph: ${err} - ${hash}`)
    }
    mermaid.initialize({
        startOnLoad: false,
        theme: cfg.theme,
        securityLevel: 'loose',
    });
}

function subscribeToEvents() {
    let eventSource = new EventSource("/api/v1/events");

    eventSource.onopen = (() => {
        console.debug("eventsource connection open");
    })

    eventSource.onerror = (ev => {
        if (ev.target.readyState === 0) {
            console.debug("reconnecting to eventsource");
        } else {
            console.error("eventsource error", ev);
        }
    })

    eventSource.onmessage = (ev => {
        let obj = JSON.parse(ev.data);
        switch (true) {
            case obj.forceReload:
                eventSource.close()
                window.location.reload()
                break
            case obj.reloadConfig:
                getRevealConfig().then(cfg => {
                    Reveal.configure(cfg)
                })
                break
            default:
                switch (true) {
                    case obj.file.endsWith(".css"):
                        let cssLink = document.querySelector(`link[rel=stylesheet][id="${obj.fileNameHash}"]`);
                        cssLink.href = `${obj.file}?ts=${obj.ts}`
                        break
                    default:
                        let elem = document.getElementById(obj.fileNameHash);
                        if (elem !== null) {
                            elem.src = `${obj.file}?ts=${obj.ts}`
                        }
                }
        }
    })
}