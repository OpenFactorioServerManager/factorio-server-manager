const path = require('path');
const MiniCssExtractPlugin = require("mini-css-extract-plugin");
var OptimizeCssAssetsPlugin = require('optimize-css-assets-webpack-plugin');
/**
 * TODO remove when webpack fixed this error:
 * Links to this error:
 * https://github.com/webpack/webpack/issues/7300
 * https://github.com/JeffreyWay/laravel-mix/pull/1495
 * https://github.com/webpack-contrib/mini-css-extract-plugin/issues/151
 * and more...
 *
 * As far as i know, this will be fixed in webpack 5
 * ~knoxfighter
 */
const FixStyleOnlyEntriesPlugin = require("webpack-fix-style-only-entries");

module.exports = (env, argv) => {
    const isProduction = argv.mode == 'production';

    return {
        entry: {
            bundle: './ui/index.js',
            style: './ui/index.scss'
        },
        output: {
            filename: '[name].js',
            path: path.resolve(__dirname, 'app'),
            publicPath: ""
        },
        resolve: {
            alias: {
                Utilities: path.resolve('ui/js/')
            },
            extensions: ['.js', '.json', '.jsx']
        },
        devtool: (isProduction) ? "none" : "source-map",
        module: {
            rules: [
                {
                    test: /\.jsx?$/,
                    exclude: /node_modules/,
                    use: {
                        loader: 'babel-loader',
                        options: {
                            presets: [
                                '@babel/preset-env',
                                [
                                    '@babel/preset-react', {
                                        development: !isProduction
                                    }
                                ]
                            ]
                        }
                    }
                },
                {
                    test: /\.scss$/,
                    use: [
                        MiniCssExtractPlugin.loader,
                        {
                            loader: "css-loader",
                            options: {
                                "sourceMap": !isProduction,
                            }
                        },
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
                                    if (!/node_modules/.test(loader_path)) {
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
        stats: {
            children: false
        },
        plugins: [
            new FixStyleOnlyEntriesPlugin(),
            new MiniCssExtractPlugin({
                filename: "[name].css"
            }),
            (isProduction) ?
                new OptimizeCssAssetsPlugin({
                    cssProcessorPluginOptions: {
                        preset: ['default', { discardComments: { removeAll: true } }],
                    },
                }) : () => {}
        ]
    }
}
