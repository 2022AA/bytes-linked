<template>
  <div class="container">
    <div class="order-list">
      <div class="img-flex">
        <el-image class="list-img" :src="avator"></el-image>
      </div>
      <div class="order-desc">
        <p>天宫数字空间站</p>
        <p>元气艺术</p>
        <p>¥19.9</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { TabsPaneContext } from 'element-plus';
import avator from '@/assets/img/avatar.jpeg';

const handleClick = (tab: TabsPaneContext, event: Event) => {
  console.log(tab, event)
}
</script>

<style lang="scss" scoped>
.container {
  width: 100vw;
  height: 100vh;
  background-color: #b5d4ef;
  padding: 1rem 1.5rem;

  .order-list {
    position: relative;
    display: flex;
    align-items: center;
    height: 5.5rem;
    padding-left: 0.6rem;
    background-color: #5ea6e1;
    border-radius: 6px;

    .img-flex {
      height: 4.5rem;
      width: 4.5rem;
      margin: 0 1rem 0 0;

      .list-img {
        height: 100%;
        width: 100%;
        border-radius: 6px;
      }
    }

    .order-desc {
      color: #fff;

      &>:nth-child(2) {
        font-size: 0.9rem;
      }

      &>:nth-child(3) {
        font-size: 0.9rem;
      }
    }
  }
}
</style>
