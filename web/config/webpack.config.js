const path = require('path')
const { CleanWebpackPlugin } = require('clean-webpack-plugin')
const HTMLWebpackPlugin = require('html-webpack-plugin')
const MiniCssExtractPlugin = require('mini-css-extract-plugin')

module.exports = {
	entry: './src/index.ts',
	devtool: 'source-map',
	module: {
		rules: [
			{
				test: /\.tsx?$/,
				use: 'ts-loader',
				exclude: /node_modules/,
			},
			{
				test: /\.(js)$/,
				exclude: /node_modules/,
				use: ['babel-loader', "eslint-loader"]
			},
			{
				test: /\.(jpg|png)$/,
				use: {
				  loader: 'url-loader',
				  options: {
					  esModule: false,
				  }
				},
			},
			{
				test: /\.scss$/,
				use: [
					MiniCssExtractPlugin.loader,
					'css-loader',
					'sass-loader',
				]
			},
		]
	},
	resolve: {
		extensions: ['*', '.js', ".ts"]
	},
	plugins: [
		new CleanWebpackPlugin(),
		new HTMLWebpackPlugin({
			template: "assets/html/index.html",
		}),
	],
	output: {
	  path: path.resolve(__dirname, "../", 'dist'),
	  publicPath: '/',
	  filename: 'bundle.js'
	},
	devServer: {
	  contentBase: './dist'
	}
  };