const path = require('path');

module.exports = {
    entry: './ui/index.js',
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
            }
        ]
    },
    performance: {
        hints: false
    }
}
