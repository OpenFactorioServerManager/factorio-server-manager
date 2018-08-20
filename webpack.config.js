const path = require('path');
const MiniCssExtractPlugin = require("mini-css-extract-plugin");

module.exports = {
    entry: {
        // js: './ui/index.js',
        sass: './ui/index.scss'
    },
    output: {
        filename: 'bundle.js',
        path: path.resolve(__dirname, 'app')
    },
    resolve: {
        alias: {
            Utilities: path.resolve(__dirname, 'ui/js/')
        },
        extensions: ['.wasm', '.mjs', '.js', '.json', '.jsx']
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
                    "sass-loader"
                ]
            },
            {
                test: /(\.(png|jpe?g|gif)$|^((?!font).)*\.svg$)/,
                loaders: [
                    {
                        loader: "file-loader",
                        options: {
                            name: path => {
                                console.log("es ist ein BILD!");
                                console.log(path);
                                if(!/node_modules|bower_components/.test(path)) {
                                    return "app/images/[name].[ext]?[hash]";
                                }

                                return (
                                    "app/images/vendor/" +
                                    path.replace(/\\/g, "/")
                                        .replace(
                                            /((.*(node_modules|bower_components))|images|image|img|assets)\//g,
                                            ''
                                        ) +
                                    '?[hash]'
                                );
                            },
                            publicPath: 'app/'
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
                            name: path => {
                                console.log("testausgabe");
                                console.log(path);
                                if (!/node_modules|bower_components/.test(path)) {
                                    return 'app/fonts/[name].[ext]?[hash]';
                                }

                                return (
                                    'app/fonts' +
                                    '/vendor/' +
                                    path
                                        .replace(/\\/g, '/')
                                        .replace(
                                            /((.*(node_modules|bower_components))|fonts|font|assets)\//g,
                                            ''
                                        ) +
                                    '?[hash]'
                                );
                            },
                            publicPath: 'app/'
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
        new MiniCssExtractPlugin({
            filename: "bundle.css"
        })
    ]
}
