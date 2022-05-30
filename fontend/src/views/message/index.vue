<template>
  <div class="container">
    <template v-for="(item, index) in messageList" :key="index">
      <div class="list-flex">
        <div class="list-title">
          <div>{{item.title}}</div>
          <div>{{item.date}}</div>
        </div>
        <div class="list-content">{{item.content}}</div>
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">

  const messageList = reactive([
    { title: '系统消息', date: '2022-03-06 17:33', content: '实名认证已启用人脸识别，如有实名认证未通过的，请调整认证环境再次尝试。' },
    { title: '系统消息', date: '2022-05-16 17:33', content: '实名认证已启用人脸识别，如有实名认证未通过的，请调整认证环境再次尝试。' },
    { title: '系统消息', date: '2022-05-16 17:33', content: '实名认证已启用人脸识别，如有实名认证未通过的，请调整认证环境再次尝试。实名认证已启用人脸识别，如有实名认证未通过的，请调整认证环境再次尝试。' },
  ])
</script>

<style lang="scss" scoped>
.container {
  width: 100vw;
  height: 100vh;
  background-color: #b5d4ef;
  padding: 1rem;

  .list-flex {
    font-size: 0.9rem;
    width: 100%;
    color: #fff;
    padding: 0.6rem;
    margin-bottom: 0.8rem;
    background-color: #5ea6e1;
    border-radius: 6px;

    .list-title {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 1rem;
    }

    .list-content {
      text-align: justify;
    }
  }
}
</style>
