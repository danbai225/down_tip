module.exports = {
    publicPath: '/socket',
    // 将构建好的文件输出到哪里
    outputDir: '../dist',
    // 放置静态资源的地方 (js/css/img/font/...)
    assetsDir: 'static',
    lintOnSave: true,
    // 使用带有浏览器内编译器的完整构建版本
    // 查阅 https://cn.vuejs.org/v2/guide/installation.html#运行时-编译器-vs-只包含运行时
    runtimeCompiler: false,
    // babel-loader 默认会跳过 node_modules 依赖。
    // 通过这个选项可以显式转译一个依赖。
    transpileDependencies: [/* string or regex */],
    // 是否为生产环境构建生成 source map？
    productionSourceMap: true,
    // 调整内部的 webpack 配置。
    // 查阅 https://github.com/vuejs/vue-docs-zh-cn/blob/master/vue-cli/webpack.md
    chainWebpack: () => {},
    configureWebpack: () => {},

    // 在生产环境下为 Babel 和 TypeScript 使用 `thread-loader`
    // 在多核机器下会默认开启。
    parallel: require('os').cpus().length > 1,
  
    // PWA 插件的选项。
    // 查阅 https://github.com/vuejs/vue-docs-zh-cn/blob/master/vue-cli-plugin-pwa/README.md
    pwa: {},
  
    // 配置 webpack-dev-server 行为。
    devServer: {
      open: process.platform === 'darwin',
      host: '0.0.0.0',
      port: 8080,
      https: false,
      hotOnly: false,
      // 查阅 https://github.com/vuejs/vue-docs-zh-cn/blob/master/vue-cli/cli-service.md#配置代理
      proxy: null, // string | Object
    }
  }