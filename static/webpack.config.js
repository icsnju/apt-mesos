var path = require('path');
var ExtractTextPlugin = require("extract-text-webpack-plugin");
module.exports = {
    entry: './src/main',
    output: {
        path: path.join(__dirname, './dist'),
        filename: 'build.js',
        publicPath: '/dist/'
    },

    module: {
        // 加载器
        loaders: [
            { test: /\.vue$/, loader: 'vue' },
            { test: /\.js$/, loader: 'babel', exclude: /node_modules/ },
            { test: /\.css$/, loader: 'style!css!autoprefixer'},
            { test: /\.scss$/, loader: ExtractTextPlugin.extract("style-loader", 'css-loader!sass-loader')},
            { test: /\.css$/, loader: ExtractTextPlugin.extract("style-loader", "css-loader")},
            { test: /\.(png|jpg|gif)$/, loader: 'url-loader'},
            { test: /\.(html|tpl)$/, loader: 'html-loader' },
            { test: /\.json$/, loader: 'json' },
        ]
    },
    vue: {
        css: ExtractTextPlugin.extract("css"),
        sass: ExtractTextPlugin.extract("css!sass-loader")
    },
    babel: {
        presets: ['es2015'],
        plugins: ['transform-runtime']
    },
    resolve: {
        extensions: ['', '.js', '.vue'],
        alias: {
            filter: path.join(__dirname, './src/filters'),
            components: path.join(__dirname, './src/components')
        }
    },
    plugins: [
        new ExtractTextPlugin("styles.css")
    ],
    devtool: '#source-map'
};