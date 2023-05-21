const webpack = require('webpack');

const { defineConfig } = require('@vue/cli-service')
module.exports = defineConfig({
  transpileDependencies: true,
  configureWebpack: {
        plugins: [
            new webpack.DefinePlugin({
                'process.env': {
                    API_URL: JSON.stringify(process.env.API_URL)
                }
            })
        ]
    }
})
