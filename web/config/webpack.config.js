const path = require('path')
const { CleanWebpackPlugin } = require('clean-webpack-plugin')
const HTMLWebpackPlugin = require('html-webpack-plugin')
const ForkTsCheckerWebpackPlugin = require('fork-ts-checker-webpack-plugin');
const VueLoaderPlugin = require('vue-loader/lib/plugin')

const rel = (...objs) => path.resolve(__dirname, "..", ...objs)

module.exports = {
    entry: {
        index: rel("src", "index", "index.ts"),
        lost: rel("src", "lost", "index.ts"),
        get: rel("src", "get", "index.ts"),
        put: rel("src", "put", "index.ts"),
    },
    output: {
        // would this work instead?
        path: rel("dist"),
    },
    devtool: 'source-map',
    module: {
        rules: [
            {
                test: /\.vue$/,
                exclude: /node_modules/,
                use: 'vue-loader',
            },
            {
                test: /\.ts$/,
                exclude: /node_modules/,
                use: 'ts-loader',
            },
            {
                test: /\.js$/,
                exclude: /node_modules/,
                use: ['babel-loader', "eslint-loader"]
            },
            {
                test: /\.(jpg|png)$/,
                exclude: /node_modules/,
                use: {
                    loader: 'url-loader',
                    options: {
                        esModule: false,
                    }
                },
            },
            {
              test: /\.handlebars$/,
              exclude: /node_modules/,
              use: 'handlebars-loader',
            },
            // {
            // 	test: /\.scss$/,
            // 	use: [
            // 		MiniCssExtractPlugin.loader,
            // 		'css-loader',
            // 		'sass-loader',
            // 	]
            // },
        ]
    },
    resolve: {
        extensions: ['.js', ".ts"]
    },
    plugins: [
        new CleanWebpackPlugin(),
        new HTMLWebpackPlugin({
            template: rel("src", "index", "index.handlebars"),
            inject: true,
            chunks: ["index"],
            filename: "index.html"
        }),
        new HTMLWebpackPlugin({
            template: rel("src", "lost", "index.handlebars"),
            inject: true,
            chunks: ["lost"],
            filename: "lost.html"
        }),
        new HTMLWebpackPlugin({
            template: rel("src", "get", "index.handlebars"),
            inject: true,
            chunks: ["get"],
            filename: "get.html"
        }),
        new HTMLWebpackPlugin({
            template: rel("src", "put", "index.handlebars"),
            inject: true,
            chunks: ["put"],
            filename: "put.html"
        }),
        new ForkTsCheckerWebpackPlugin({eslint: {
            // required - same as command
            // `eslint ./src/**/*.{ts,tsx,js,jsx} --ext .ts,.tsx,.js,.jsx`
            files: './src/**/*.{ts,js}'
        }}),
        new VueLoaderPlugin(),
    ],
};