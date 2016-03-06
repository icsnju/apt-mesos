function SmoothlyMenu() {
    $("body").hasClass("mini-navbar") ? $("body").hasClass("fixed-sidebar") ? ($("#side-menu").hide(), setTimeout(function() {
        $("#side-menu").fadeIn(500)
    },
    300)) : $("#side-menu").removeAttr("style") : ($("#side-menu").hide(), setTimeout(function() {
        $("#side-menu").fadeIn(500)
    },
    100))
}

$(document).ready(function() {
    $(".navbar-minimalize").click(function() {
        $("body").toggleClass("mini-navbar"),
        SmoothlyMenu()
    })
})