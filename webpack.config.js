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
                    "sass-loader"
                ]
            },
            {
                test: /(\.(png|jpe?g|gif)$|^((?!font).)*\.svg$)/,
                loaders: [
                    {
                        loader: "file-loader",
                        options: {
                            name: loader_path => {
                                if(!/node_modules|bower_components/.test(loader_path)) {
                                    return "app/images/[name].[ext]?[hash]";
                                }

                                return (
                                    "app/images/vendor/" +
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
                                if (!/node_modules|bower_components/.test(loader_path)) {
                                    return 'app/fonts/[name].[ext]?[hash]';
                                }

                                return (
                                    'app/fonts/vendor/' +
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
        new MiniCssExtractPlugin({
            filename: "bundle.css"
        })
    ]
}
