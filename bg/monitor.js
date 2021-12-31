console.log("I'm alive")

let listeningTabs = []
let options = defaultOptions

chrome.storage.local.get(defaultOptions,
    o => options = o)


chrome.storage.onChanged.addListener((changes, area) => {
    if(area === "local") {
	const optionKeys = Object.keys(options)
	for(key of Object.keys(changes)) {
	    if(optionKeys.indexOf(key) >= 0) {
		options[key] = changes[key].newValue
	    }
	}
    }
})

chrome.browserAction.onClicked.addListener((tab) => {
        toggleTab(tab.id)
})

napp = browser.runtime.connectNative("native_inserter")
napp.onDisconnect.addListener((p) => {
        if (p.error) {
                console.error(`Disconnected native app due to an error: ${p.error.message}`)
        }
})
napp.onMessage.addListener((msg) => {
        console.log("Received: " + msg.body)
        const pasteTarget = document.querySelector("#paste-target")
        pasteTarget.innerText = msg.body
        const content = pasteTarget.innerText
        listeningTabs.forEach(id => notifyForeground(id, content))
})

function toggleTab(id) {
        const index = listeningTabs.indexOf(id)
        if(index >= 0) {
                uninject(id)
                listeningTabs.splice(index, 1)
                chrome.browserAction.setBadgeText({ text: "", tabId: id })
        } else {
                chrome.tabs.executeScript({file: "/fg/insert.js"})
                listeningTabs.push(id)
                chrome.browserAction.setBadgeBackgroundColor({ color: "green", tabId: id })
                chrome.browserAction.setBadgeText({ text: "ON", tabId: id })
        }
}

function notifyForeground(id, text) {
        chrome.tabs.sendMessage(id, { action: "insert", text, options })
}

function uninject(id) {
    chrome.tabs.sendMessage(id, { action: "uninject" })
}
