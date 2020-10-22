const { merge } = require('webpack-merge');
const common = require("./webpack.config.js")
const Dotenv = require('dotenv-webpack');
const { DefinePlugin } = require("webpack")

module.exports = merge(common, {
	mode: "production",
	output: {
		filename: "[name].[chunkhash].js"
	},
	plugins: [
		new DefinePlugin({
			'process.env': {
				'NODE_ENV': JSON.stringify('production')
			}
		}),
		new Dotenv({
			path: "./.env.prod"
		}),
		// new MiniCssExtractPlugin({
		// 	filename: '[name].[hash].css',
		// 	chunkFilename: '[name].[hash].css'
		// })
	],

	// https://webpack.js.org/configuration/optimization/
	// optimization: {
	// 	minimize: true,
	// 	minimizer: []
	// }
})