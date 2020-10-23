const path = require('path')
const { CleanWebpackPlugin } = require('clean-webpack-plugin')
const HTMLWebpackPlugin = require('html-webpack-plugin')
const ForkTsCheckerWebpackPlugin = require('fork-ts-checker-webpack-plugin');
const VueLoaderPlugin = require('vue-loader/lib/plugin')
const MiniCssExtractPlugin = require('mini-css-extract-plugin')

const rel = (...objs) => path.resolve(__dirname, "..", ...objs)

module.exports = {
  entry: {
      index: rel("src", "index", "index.ts"),
      // index: [
        // rel("src", "index", "index.ts"),
        // rel("src", "index", "index.css"),
      // ],
      // lost: [
      //   rel("src", "lost", "index.ts"),
      //   rel("src", "lost", "index.css"),
      // ],
      // get: [
      //   rel("src", "get", "index.ts"),
      //   rel("src", "get", "index.css"),
      // ],
      // put: [
      //   rel("src", "put", "index.ts"),
      //   rel("src", "put", "index.css"),
      // ],
  },
    output: {
        // would this work instead?
        path: rel("dist"),
        publicPath: '/',
    },
    devtool: 'source-map',
    module: {
        rules: [
        {
          test: /\.ts$/,
          exclude: /node_modules/,
          loader: 'ts-loader',
          options: {
            appendTsSuffixTo: [/\.vue$/]
          }
        },
        {
          test: /\.vue$/,
          use: 'vue-loader',
        },
        // It is common to have exclude: /node_modules/ for JS transpilation
        // rules (e.g. babel-loader) that apply to .js files. Due to the
        // inference change of v15, if you import a Vue SFC inside node_modules,
        // its <script> part will be excluded from transpilation as well.
        // In order to ensure JS transpilation is applied to Vue SFCs in
        // node_modules, you need to whitelist them by using an exclude function
        // instead:
        // https://vue-loader.vuejs.org/guide/pre-processors.html#babel
        {
          test: /\.js$/,
          exclude: file => (
            /node_modules/.test(file) && ~/\.vue\.js/.test(file)
          ),
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
        {
          test: /\.css$/,
          use: ['css-loader'],
        },
        {
          test: /\.scss$/,
          use: [
            {
              loader: MiniCssExtractPlugin.loader,
              options: {
                esModule: false,
              },
            },
            'css-loader',
            {
              loader: 'sass-loader',
              options: {
                additionalData: '@import \'~assets/scss/_base.scss\';',
                // sassOptions: {
                //   indentedSyntax: true
                // }
              }
            }
          ],
        },
        ]
    },
    resolve: {
        extensions: ['.js', ".ts", ".vue"],
        alias: {
          // ESM for bundlers: intended for use with modern bundlers like
          // webpack 2 or Rollup. ESM format is designed to be statically
          // analyzable so the bundlers can take advantage of that to perform
          // “tree-shaking” and eliminate unused code from your final bundle.
          // https://vuejs.org/v2/guide/installation.html#Explanation-of-Different-Builds
          'vue$': 'vue/dist/vue.esm.js',
          'assets': rel('assets'),
        }
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