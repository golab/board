const observer = new MutationObserver(function(mutations) {
    let once = 0;
    mutations.forEach(function(mutation) {
        if (mutation.target.getAttribute("class") == "Goban") {
            if (once == 0) {
                once++;
                main();
            }
        }
    });
});

// Start observing the document
observer.observe(document.body, {
    childList: true,
    subtree: true
});

function main() {
    if (document.getElementById("golab") != null) {
        return;
    }
    
    let dock = document.getElementsByClassName("Dock")[0];
    if (dock == null) {
        return;
    }
    
    let tooltip_container = document.createElement("div");
    tooltip_container.setAttribute("class", "TooltipContainer");
    tooltip_container.id = "golab";
    
    let disabled = document.createElement("div");
    disabled.setAttribute("clasS", "Tooltip disabled");
    
    tooltip_container.appendChild(disabled);
    
    let p = document.createElement("p");
    p.setAttribute("class", "title");
    p.innerHTML = "Upload to Go Lab";
    
    let div = document.createElement("div");
    let anchor = document.createElement("a");
    anchor.href = "https://golab.gg/ext/upload?url=" + window.location.href;
    anchor.target = "_blank";
    
    let icon = document.createElement("i");
    icon.setAttribute("class", "fa fa-upload");
    
    anchor.appendChild(icon);
    anchor.innerHTML += "Upload to Go Lab";
    
    div.appendChild(anchor);
    tooltip_container.appendChild(div);
    
    dock.appendChild(tooltip_container);
}
