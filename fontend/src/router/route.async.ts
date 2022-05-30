// 需要鉴权的业务路由
import { RouteRecordRaw } from 'vue-router';

const asyncRoutes: Array<RouteRecordRaw> = [
  {
    path: '/',
    name: 'login',
    meta: {
      title: '',
      icon: '',
    },
    component: () => import('@/views/login/index.vue'),
  },
  {
    path: '/home',
    name: 'home',
    meta: {
      title: '',
      icon: '',
    },
    component: () => import('@/views/home/index.vue'),
  },
  {
    path: '/collect',
    name: 'collect',
    meta: {
      title: '',
      icon: '',
    },
    component: () => import('@/views/collect/index.vue'),
  },
  {
    path: '/trade',
    name: 'trade',
    meta: {
      title: '',
      icon: '',
    },
    component: () => import('@/views/trade/index.vue'),
  },
  {
    path: '/upload',
    name: 'upload',
    meta: {
      title: '',
      icon: '',
    },
    component: () => import('@/views/upload/index.vue'),
  },
  {
    path: '/center',
    name: 'center',
    meta: {
      title: '',
      icon: '',
    },
    component: () => import('@/views/center/index.vue'),
  },
  {
    path: '/order',
    name: 'order',
    meta: {
      title: '',
      icon: '',
    },
    component: () => import('@/views/order/index.vue'),
  },
  {
    path: '/record',
    name: 'record',
    meta: {
      title: '',
      icon: '',
    },
    component: () => import('@/views/record/index.vue'),
  },
  {
    path: '/message',
    name: 'message',
    meta: {
      title: '',
      icon: '',
    },
    component: () => import('@/views/message/index.vue'),
  },
  {
    path: '/hisCollect',
    name: 'hisCollect',
    meta: {
      title: '',
      icon: '',
    },
    component: () => import('@/views/hisCollect/index.vue'),
  },
  {
    path: '/process',
    name: 'process',
    meta: {
      title: 'Template configuration process',
      icon: '',
    },
    component: () => import('@/views/example/MarkdownPage.vue'),
  },
];

export default asyncRoutes;
