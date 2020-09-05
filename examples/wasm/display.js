function init() {
    const canvas = document.getElementById('screen');
    let ctx = canvas.getContext('2d');
    document.ctx = ctx;
    canvas.width = 640;
    canvas.height = 320;
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

function repaint() {
    for (let i = 0; i < 64 * 32; i++) {
        const x = i % 64;
        const y = Math.floor(i / 64);
        if (document.vram[i] === 0) {
            document.ctx.fillStyle = 'black';
        } else {
            document.ctx.fillStyle = 'white';
        }
        ctx.fillRect(x * 10, y * 10, 10, 10)
    }
}
