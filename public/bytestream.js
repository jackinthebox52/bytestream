var divCheckingInterval = setInterval(function(){
    // Find it with a selector
    var elem = document.querySelector(".alert")
    if(elem){
        console.log("Found!");
        clearInterval(divCheckingInterval);
        reloadPlayer(elem);
    }
}, 500);

function reloadPlayer(elem){
    var player = document.querySelector(".byteplayer");
    if(player){
       player.src = player.src
    }
    return;
}
