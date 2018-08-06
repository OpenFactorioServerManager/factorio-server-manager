const mix = require('laravel-mix');

mix.setPublicPath("app");

mix.react('ui/index.js', 'bundle.js')
   .sass('ui/index.scss', 'bundle.css')
   // .copyDirectory('node_modules/bootstrap-fileinput/img', 'app/img');

if (!mix.inProduction()) {
    mix.webpackConfig({
        devtool: 'source-map'
    })
        .sourceMaps()
}
