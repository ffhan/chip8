let context = new AudioContext()
let o = context.createOscillator()
let g = context.createGain()

function playMusic() {
    context = new AudioContext()
    o = context.createOscillator()
    g = context.createGain()
    o.connect(g)
    g.connect(context.destination)
    o.start(0)
}

function stopMusic() {
    g.gain.exponentialRampToValueAtTime(
        0.00001, context.currentTime + 0.04
    )
}
