function init() {
    const canvas = document.getElementById('screen');
    canvas.width = 640;
    canvas.height = 320;
    let ctx = canvas.getContext('2d');
    ctx.fillRect(0, 0, 640, 320);
}

async function showHelp(name) {
    const response = await fetch(`help/${name}.html`)
    if (!response.ok) {
        clearHelp();
        return;
    }
    const responseHtml = await response.text();

    $('#gameHelp').fadeOut(500, function () {
        $(this).html(responseHtml).fadeIn(500);
    });
    $('#gameHelp').prop('hidden', false);
}

function clearHelp() {
    $('#gameHelp').fadeOut(500, function () {
        $('#gameHelp').prop('hidden', true);
    });
}
