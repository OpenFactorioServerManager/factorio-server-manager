const mix = require('laravel-mix');

mix.react('ui/index.js', 'app/bundle.js')
    .sass('ui/index.scss', 'app/bundle.css');

if (!mix.inProduction()) {
    mix.webpackConfig({
        devtool: 'source-map'
    })
        .sourceMaps()
}
