function init() {
    const canvas = document.getElementById('screen');
    canvas.width = 640;
    canvas.height = 320;
    let ctx = canvas.getContext('2d');
    ctx.fillRect(0, 0, 640, 320);
}