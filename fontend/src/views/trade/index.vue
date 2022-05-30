<template>
  <div class="container">
    <div class="list-title">最新交易</div>
    <div>
      <template v-for="(item, index) in tradeList" :key="index">
        <div class="trade-list-flex">
          <div class="list-desc">
            <div class="desc-top">
              <el-avatar class="list-avatar" :src="item.avator"></el-avatar>
              <p>{{item.author}}完成了一笔交易</p>
            </div>
            <p>作品名称：{{item.title}}</p>
            <p>交易时间：{{item.tradeTime}}</p>
          </div>
          <div class="img-flex">
            <el-image class="list-img" :src="item.img" />
          </div>
        </div>
      </template>
    </div>
    <div class="list-title">热门作品</div>
    <div>
      <el-row justify="space-between" :gutter="24">
        <template v-for="(item, index) in modelList" :key="index">
          <el-col :span="12">
            <div class="model-list-flex">
              <div class="list-img">
                <el-image :src="item.img" />
              </div>
              <div class="list-desc">
                <p>{{item.title}}</p>
                <p>{{item.author}}</p>
                <span class="price">¥{{item.price}}</span>
              </div>
            </div>
          </el-col>
        </template>
      </el-row>
    </div>
  </div>
</template>

<script setup lang="ts">
  import avator from '@/assets/img/avatar.jpeg';
  import demo from '@/assets/img/img-demo.png'

  const { t } = useI18n();
  // ElMessage.info({ message: 'welcome', duration: 1 });
  // ElNotification({
  //   title: 'Issue',
  //   message: 'If you encounter problems in using the template, please raise them in the issue',
  //   duration: 0,
  // })

  const tradeList = reactive([
    { avator: avator, img: demo, title: '作品一', author: '2022AA', tradeTime: '2022-05-16 23:59:59' },
    { avator: avator, img: demo, title: '作品一', author: '2022AA', tradeTime: '2022-05-16 23:59:59' },
  ])

  const modelList = reactive([
    { img: demo, title: '作品一', limit: 1000, author: '2022AA', price: 100 },
    { img: demo, title: '作品一', limit: 1000, author: '2022AA', price: 100 },
    { img: demo, title: '作品一', limit: 1000, author: '2022AA', price: 100 },
    { img: demo, title: '作品一', limit: 1000, author: '2022AA', price: 100 },
    { img: demo, title: '作品一', limit: 1000, author: '2022AA', price: 100 },
    { img: demo, title: '作品一', limit: 1000, author: '2022AA', price: 100 },
    { img: demo, title: '作品一', limit: 1000, author: '2022AA', price: 100 },
    { img: demo, title: '作品一', limit: 1000, author: '2022AA', price: 100 },
    { img: demo, title: '作品一', limit: 1000, author: '2022AA', price: 100 },
  ]);
</script>

<style lang="scss" scoped>
  .container {
    position: fixed;
    width: 100vw;
    height: 100vh;
    left: 0;
    bottom: 4rem;
    overflow: auto;
    background-color: #b5d4ef;
    padding: 4rem 1rem 0;

    .list-title {
      color: #fff;
      margin: 2rem 0 1rem;
    }

    .trade-list-flex {
      display: flex;
      justify-content: space-between;
      align-items: center;
      width: 100%;
      border-radius: 6px;
      background-color: #5ea6e1;
      margin-bottom: 1rem;
      padding: 0.8rem;

      .list-desc {
        position: relative;
        color: #fff;

        .desc-top {
          display: flex;
          align-items: center;
        }

        p{
          margin-bottom: 0.3rem;
        }

        &>:first-child {
          font-size: 0.9rem;
          margin-bottom: 0.6rem;

          p {
            margin-left: 0.6rem;
          }
        }

        &>:last-child {
          font-size: 0.9rem;
        }
      }

      .img-flex {
        width: 6rem;
        margin-left: 1rem;

        list-img {
          width: 100%;
          height: auto;
          border-radius: 10px;
        }
      }
    }

    .model-list-flex {
      width: 100%;
      border-radius: 6px;
      background-color: #5ea6e1;
      margin-bottom: 1rem;
      padding: 0.8rem;

      .list-img {
        width: 100%;
        height: auto;
      }

      .list-desc {
        position: relative;
        color: #fff;

        p{
          margin-bottom: 0.3rem;
        }

        &>:first-child {
          font-size: 1.1rem;
        }

        &>:last-child {
          font-size: 0.9rem;
        }

        .price {
          position: absolute;
          right: 0.5rem;
          top: 50%;
          transform: translateY(-50%);
          font-size: 1.2rem;
          font-weight: bold;
        }
      }
    }
  }
</style>
