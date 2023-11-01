const path = require('path');
const MiniCssExtractPlugin = require("mini-css-extract-plugin");
const TerserPlugin = require('terser-webpack-plugin');

module.exports = (env, argv) => {
    const isProduction = argv.mode === 'production';

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
        devtool: isProduction ? false : "source-map",
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
                        {
                            loader: "sass-loader",
                            options: {
                                // always make sourceMap. resolver-url-loader is needing it
                                "sourceMap": true,
                            }
                        },
                        {
                            loader: 'postcss-loader',
                            options: {
                                postcssOptions: {
                                    plugins: [
                                        require('tailwindcss'),
                                        require('autoprefixer'),
                                    ],
                                }
                            },
                        }
                    ]
                },
                {
                    test: /(\.(png|jpe?g|gif)$|^((?!font).)*\.svg$)/,
                    use: [
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
                    use: [
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
        plugins: [new MiniCssExtractPlugin()],
        optimization: {
            minimize: isProduction,
            minimizer: [
                new MiniCssExtractPlugin(
                    {
                        filename: "[name].css"
                    }
                ),
                new TerserPlugin()
            ],
        }
    }
}
