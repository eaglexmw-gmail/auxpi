'use strict'
const path = require('path')
const utils = require('./utils')
const webpack = require('webpack')
const config = require('../config')
const merge = require('webpack-merge')
const baseWebpackConfig = require('./webpack.base.conf')
const CopyWebpackPlugin = require('copy-webpack-plugin')
const HtmlWebpackPlugin = require('html-webpack-plugin')
const MiniCssExtractPlugin = require('mini-css-extract-plugin')
const OptimizeCSSPlugin = require('optimize-css-assets-webpack-plugin')
const UglifyJsPlugin = require('uglifyjs-webpack-plugin')

const env = require('../config/prod.env')

const webpackConfig = merge(baseWebpackConfig, {
  mode: 'production',
  module: {
    /*
    rules: utils.styleLoaders({
      sourceMap: config.build.productionSourceMap,
      extract: true,
      usePostCSS: true
    })
    */
    rules: [
      {
        test: /\.css$/,
        //include: [path.resolve(__dirname, 'src')],   // 限制打包范围，提高打包速度
        //exclude: /node_modules/,                     // 排除node_modules文件夹
        use: [
            // {    // 当配置MiniCssExtractPlugin.loader后，此项就无需配置，原因看各自作用
            //     loader: "style-loader"  // 将处理结束的css代码存储在js中，运行时嵌入`<style>`后挂载到html页面上
            // },
            {
                loader: MiniCssExtractPlugin.loader  // 将处理后的CSS代码提取为独立的CSS文件，可以只在生产环境中配置，但我喜欢保持开发环境与生产环境尽量一致
            },
            {
                loader: "css-loader"    // CSS加载器，使webpack可以识别css文件
            }
        ]
      },
      {
        test: /\.scss$/,
        //include: [path.resolve(__dirname, 'src')],   // 限制打包范围，提高打包速度
        //exclude: /node_modules/,                     // 排除node_modules文件夹
        use: [
            // {    // 当配置MinCssExtractPlugin.loader后，此项就无需配置，原因看各自作用
            //     loader: "style-loader"  // 将处理结束的css代码存储在js中，运行时嵌入`<style>`后挂载到html页面上
            // },
            {
                loader: MiniCssExtractPlugin.loader,  // 将处理后的CSS代码提取为独立的CSS文件，可以只在生产环境中配置，但我喜欢保持开发环境与生产环境尽量一致
            },
            {
                loader: "css-loader"    // CSS加载器，使webpack可以识别css文件
            },
            {   
                loader: "postcss-loader"    //承载autoprefixer功能，为css添加前缀
            },
            {
                loader: "sass-loader",       // 编译sass，webpack默认使用node-sass进行编译，所以需要同时安装 sass-loader 和 node-sass
                options: {      // loader 的额外参数，配置视具体 loader 而定
                    sourceMap: true, // 要安装resolve-url-loader，当此配置项启用 sourceMap 才能正确加载 Sass 里的相对路径资源，类似background: url(../image/test.png)
                }
            }
        ]
      },
    ]
  },
  devtool: config.build.productionSourceMap ? config.build.devtool : false,
  output: {
    path: config.build.assetsRoot,
    filename: utils.assetsPath('js/[name].[chunkhash:8].js'),
    chunkFilename: utils.assetsPath('js/chunk-[chunkhash:4].[contenthash:8].js'),
  },
  optimization: {
    //minimize: false,
    splitChunks: {
      chunks: "all",
      cacheGroups: {
        libs: { // 第三方库
          name: "chunk-libs",
          test: /[\\/]node_modules[\\/]/,
          priority: 10,
          filename: 'js/[name].[chunkhash:8].js',
          chunks: "initial" // 只打包初始时依赖的第三方
        },
        elementUI: { // elementUI 单独拆包
          name: "chunk-elementUI",
          test: /[\\/]node_modules[\\/]element-ui[\\/]/,
          priority: 20, // 权重要大于 libs
          filename: 'js/[name].[chunkhash:8].js',
          chunks: "initial" // 只打包初始时依赖的第三方
        },
        default: {
          minChunks: 2, //被引用两次就提取出来
          priority: -20,
          reuseExistingChunk: true, //检查之前是否被引用过有的话就不被打包了
        },
        styles: {
          name: 'styles',
          test: /\.css$/,
          chunks: 'initial',
          enforce: true
        },
      }
    },
  },
  plugins: [
    // http://vuejs.github.io/vue-loader/en/workflow/production.html
    new webpack.DefinePlugin({
      'process.env': env
    }),
    //打包环境去掉console.log
    new UglifyJsPlugin({
      uglifyOptions: {
        compress: {
          warnings: false,
          drop_console: true,  //注释console
          drop_debugger: true, //注释debugger
          pure_funcs: ['console.log'], //移除console.log
        },
      },
    }),
    // extract css into its own file
    new MiniCssExtractPlugin({
      // Options similar to the same options in webpackOptions.output
      // both options are optional
      filename: utils.assetsPath('css/[name].[chunkhash:8].css'),
      chunkFilename: utils.assetsPath('css/chunk-[chunkhash:4].[contenthash:8].css'),
    }),
    // Compress extracted CSS. We are using this plugin so that possible
    // duplicated CSS from different components can be deduped.
    new OptimizeCSSPlugin({
      cssProcessorOptions: config.build.productionSourceMap
        ? { safe: true, map: { inline: false } }
        : { safe: true }
    }),
    // generate dist index.html with correct asset hash for caching.
    // you can customize output by editing /index.html
    // see https://github.com/ampedandwired/html-webpack-plugin
    new HtmlWebpackPlugin({
      filename: config.build.index,
      template: 'index.html',
      inject: true,
      minify: {
        removeComments: true, // 移除注释
        collapseWhitespace: true, // 折叠空白区域
        removeAttributeQuotes: true, // 移除属性的引号
        collapseBooleanAttributes: true, // 省略只有 boolean 值的属性值 例如：readonly checked
        removeRedundantAttributes:true, // 删除多余的属性
        // more options:
        // https://github.com/kangax/html-minifier#options-quick-reference
      },
      // necessary to consistently work with multiple chunks via CommonsChunkPlugin
      chunksSortMode: 'none'
    }),
    // keep module.id stable when vendor modules does not change
    new webpack.HashedModuleIdsPlugin(),
    // enable scope hoisting
    //new webpack.optimize.ModuleConcatenationPlugin(),

    // copy custom static assets
    new CopyWebpackPlugin([
      {
        from: path.resolve(__dirname, '../static'),
        to: config.build.assetsSubDirectory,
        ignore: ['.*']
      }
    ])
  ]
})

if (config.build.productionGzip) {
  const CompressionWebpackPlugin = require('compression-webpack-plugin')

  webpackConfig.plugins.push(
    new CompressionWebpackPlugin({
      filename: '[path].gz[query]',
      algorithm: 'gzip',
      test: new RegExp(
        '\\.(' +
        config.build.productionGzipExtensions.join('|') +
        ')$'
      ),
      threshold: 10240,
      minRatio: 0.8
    })
  )
}

if (config.build.bundleAnalyzerReport) {
  const BundleAnalyzerPlugin = require('webpack-bundle-analyzer').BundleAnalyzerPlugin
  webpackConfig.plugins.push(new BundleAnalyzerPlugin())
}

module.exports = webpackConfig
