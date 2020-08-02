module.exports = {
  publicPath: '/ui/',
  devServer: {
    proxy: 'http://dashboard_backend:8080',
  },
};
