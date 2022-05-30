<template>
  <div id="container">
    <el-config-provider :locale="locale">
      <router-view></router-view>
    </el-config-provider>
    <template v-if="routerName !== '/' && routerName !== '/login' && routerName !== '/upload'">
      <div class="menu-flex">
        <el-row justify="space-between" :gutter="24">
          <el-col
              class="menu-item"
              :class="routerName === '/home' ? 'active-item' : ''"
              :span="8"
              @click="onClickMenu('/home')"
          >首页</el-col>
          <el-col
              class="menu-item"
              :class="routerName === '/collect' ? 'active-item' : ''"
              :span="8"
              @click="onClickMenu('/collect')"
          >作品</el-col>
          <el-col
              class="menu-item"
              :class="routerName === '/center' ? 'active-item' : ''"
              :span="8"
              @click="onClickMenu('/center')"
          >我的</el-col>
        </el-row>
      </div>
    </template>
  </div>
</template>
<script setup lang="ts">
  import zhCn from 'element-plus/lib/locale/lang/zh-cn';
  import { useRouter } from 'vue-router';
  import { storeToRefs } from "pinia";
  import routerStore from '@/store/router';
  import userStore from '@/store/user';

  const locale = zhCn;
  const router = useRouter();

  // 解构并使Store数据具有响应式
  const { routerName } = storeToRefs(routerStore());

  // 页面初始化加载
  onMounted(() => {
    // 初始化根据当前地址栏pathname更新routerName
    routerStore().$patch({
      routerName: (window.location.hash).substring(1),
    })

    const storageUserName: string = window.localStorage.getItem('username') || '';
    const storageToken: string = window.localStorage.getItem('token') || '';
    const storageUid: string = window.localStorage.getItem('uid') || '';
    const storagePhone: string = window.localStorage.getItem('phone') || '';
    const storageInviteCode: string = window.localStorage.getItem('inviteCode') || '';
    const storageAvatarUrl: string = window.localStorage.getItem('avatarUrl') || '';
    const storageBalance: string = window.localStorage.getItem('balance') || '';

    if (storageUserName !== null && storageUserName !== '') {
      userStore().$patch({
        username: storageUserName,
        token: storageToken,
        uid: storageUid,
        phone: storagePhone,
        inviteCode: storageInviteCode,
        avatarUrl: storageAvatarUrl,
        balance: storageBalance,
      })
    } else {
      // 登录失效，跳转登录页
      router.push('/');
    }
  })

  /**
   * 切换菜单
   * @param path
   */
  const onClickMenu = (path: string) => {
    // 更新路由Store
    routerStore().$patch({
      routerName: path,
    })
    // 跳转路由页
    router.push(path);
  }

</script>
<style lang="scss" scoped>
  .menu-flex {
    position: fixed;
    bottom: 0;
    left: 0;
    width: 100vw;
    height: 3.6rem;
    background-color: #f7f7f7;
    border-top: 1px solid #e6e6e6;

    .menu-item {
      font-size: 1rem;
      height: 3.6rem;
      line-height: 3.6rem;
      text-align: center;
    }

    .active-item {
      color: #4169e1;
    }

  }

  //修改组件默认的样式
  :deep(.el-form-item__error) {
    margin-top: 0.25rem;
  }
  :deep(.el-dialog) {
    border-radius: 6px;
  }
  :deep(.el-dialog__body) {
    padding: 1rem;
  }
  :deep(.el-dialog__header) {
    padding-top: 12px;
  }
  :deep(.el-dialog__headerbtn) {
    top: 0;
    width: 52px;
    height: 52px;
  }
  :deep(.el-dialog__footer) {
    padding-bottom: 10px;
  }
  :deep(.el-tabs__content) {
    color: #6b778c;
    font-size: 1.2rem;
    font-weight: 600;
  }
</style>
