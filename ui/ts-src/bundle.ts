import chokidar from 'chokidar';

const prod = process.env.NODE_ENV === 'production';

async function bundle() {
  const build = await Bun.build({
    entrypoints: ['./ts-src/index.ts'],
    outdir: './static/dist',
    minify: prod,
    splitting: true,
    naming: {
      asset: '[name].[ext]',
    }
  });

  for (const artifact of build.outputs) {
    console.log(`Wrote ${artifact.path}`)
  }
}

if (prod) {
  await bundle();
} else {
  const watcher = chokidar.watch('./ts-src', {
    persistent: true,
  });

  watcher.on('change', bundle);
  watcher.on('add', bundle);
  watcher.on('unlink', bundle);
}

export { };
