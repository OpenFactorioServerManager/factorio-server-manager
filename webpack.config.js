const path = require('path');
const MiniCssExtractPlugin = require("mini-css-extract-plugin");

/* TODO remove when webpack fixed this error:
 * Links to this error:
 * https://github.com/webpack/webpack/issues/7300
 * https://github.com/JeffreyWay/laravel-mix/pull/1495
 * https://github.com/webpack-contrib/mini-css-extract-plugin/issues/151
 * and more...
 *
 * This will be, as far as i know, fixed in webpack 5, it is currently in development
 */
const FixStyleOnlyEntriesPlugin = require("webpack-fix-style-only-entries");

module.exports = {
    entry: {
        // js: './ui/index.js',
        sass: './ui/index.scss'
    },
    output: {
        filename: 'bundle.js',
        path: path.resolve(__dirname, 'app'),
        publicPath: ""
    },
    resolve: {
        alias: {
            Utilities: path.resolve('ui/js/')
        },
        extensions: ['.js', '.json', '.jsx']
    },
    module: {
        rules: [
            {
                test: /\.jsx?$/,
                exclude: /node_modules/,
                use: {
                    loader: 'babel-loader'
                }
            },
            {
                test: /\.scss$/,
                use: [
                    MiniCssExtractPlugin.loader,
                    "css-loader",
                    "resolve-url-loader",
                    "sass-loader?sourceMap"
                ]
            },
            {
                test: /(\.(png|jpe?g|gif)$|^((?!font).)*\.svg$)/,
                loaders: [
                    {
                        loader: "file-loader",
                        options: {
                            name: loader_path => {
                                if(!/node_modules/.test(loader_path)) {
                                    return "/images/[name].[ext]?[hash]";
                                }

                                return (
                                    "/images/vendor/" +
                                    loader_path.replace(/\\/g, "/")
                                        .replace(/((.*(node_modules))|images|image|img|assets)\//g, '') +
                                    '?[hash]'
                                );
                            },
                        }
                    }
                ]
            },
            {
                test: /(\.(woff2?|ttf|eot|otf)$|font.*\.svg$)/,
                loaders: [
                    {
                        loader: "file-loader",
                        options: {
                            name: loader_path => {
                                if (!/node_modules/.test(loader_path)) {
                                    return '/fonts/[name].[ext]?[hash]';
                                }

                                return (
                                    '/fonts/vendor/' +
                                    loader_path
                                        .replace(/\\/g, '/')
                                        .replace(/((.*(node_modules))|fonts|font|assets)\//g, '') +
                                    '?[hash]'
                                );
                            },
                        }
                    }
                ]
            }
        ]
    },
    performance: {
        hints: false
    },
    plugins: [
        new FixStyleOnlyEntriesPlugin(),
        new MiniCssExtractPlugin({
            filename: "bundle.css"
        })
    ]
}
