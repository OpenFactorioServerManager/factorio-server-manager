const mix = require('laravel-mix');

mix.react('ui/index.js', 'app/bundle.js');
    // .sass('resources/sass/app.scss', 'public/css');

if (!mix.inProduction()) {
    mix.webpackConfig({
        devtool: 'source-map'
    })
        .sourceMaps()
}
