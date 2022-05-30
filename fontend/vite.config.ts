import { defineConfig, loadEnv } from 'vite';
import { resolve } from 'path';
import presets from './presets/presets';

// https://vitejs.dev/config/
export default defineConfig((env) => {
  // env 环境变量
  const viteEnv = loadEnv(env.mode, `.env.${env.mode}`);

  return {
    // base: viteEnv.VITE_BASE,
    base: './',
    // 插件
    plugins: [presets(env)],
    // envDir: './',
    // 别名设置
    resolve: {
      alias: {
        '@': resolve(__dirname, './src'), // 把 @ 指向到 src 目录去
        api: resolve(__dirname, './src/api'),
        views: resolve(__dirname, './src/views'),
        components: resolve(__dirname, './src/components'),
        utils: resolve(__dirname, './src/utils'),
        assets: resolve(__dirname, "./src/assets"),
        store: resolve(__dirname, "./src/store"),
      },
    },
    // 服务设置
    server: {
      host: true, // host设置为true才可以使用network的形式，以ip访问项目
      port: 8080, // 端口号
      open: true, // 自动打开浏览器
      cors: true, // 跨域设置允许
      strictPort: true, // 如果端口已占用直接退出
      // 接口代理
      proxy: {
        '/api': {
          // 本地 8000 前端代码的接口 代理到 8888 的服务端口
          target: 'http://192.168.56.18:8080/',
          changeOrigin: true, // 允许跨域
          rewrite: (path) => path.replace('/api', '/api/v1'),
        },
      },
    },
    build: {
      brotliSize: false,
      // 消除打包大小超过500kb警告
      chunkSizeWarningLimit: 2000,
      // 在生产环境移除console.log
      minify: 'esbuild',
      assetsDir: 'static/assets',
      outDir: 'dist',
      // 静态资源打包到dist下的不同目录
      rollupOptions: {
        output: {
          chunkFileNames: 'static/js/[name]-[hash].js',
          entryFileNames: 'static/js/[name]-[hash].js',
          assetFileNames: 'static/[ext]/[name]-[hash].[ext]',
        },
      },
    },
    css: {
      preprocessorOptions: {
        // 全局引入了 scss 的文件
        scss: {
          additionalData: `
          @import "@/assets/styles/variables.scss";
        `,
          javascriptEnabled: true,
        },
      },
    },
  };
});
