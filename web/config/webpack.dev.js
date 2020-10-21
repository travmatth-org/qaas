const { merge } = require('webpack-merge');
const common = require("./webpack.config.js")
const { DefinePlugin } = require("webpack")
const Dotenv = require('dotenv-webpack');
const MiniCssExtractPlugin = require('mini-css-extract-plugin')

module.exports = merge(common, {
	mode: "development",
	devtool: "eval-source-map",
	plugins: [
		new DefinePlugin({
			'process.env': {
				'NODE_ENV': JSON.stringify('development')
			}
		}),
		new Dotenv({
			path: "./.env.dev"
		}),
		new MiniCssExtractPlugin({
			filename: '[name].css',
			chunkFilename: '[name].css'
		})
	],
})